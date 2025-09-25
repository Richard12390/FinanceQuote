package pipeline

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"

	util "go_quotes/internal/infra/jsonutil"
	"go_quotes/internal/quotes/application/usecase"
)

type PipelineConfig struct {
	Type                string
	Symbols             []string
	Every               time.Duration
	BatchSize           int
	BatchesPerTick      int
	MaxWorkers          int
	MaxConcurrency      int
	OutputDir           string
	WriteJson           bool
	WriteJsonEveryTicks int
}

type BatchFetcher interface {
	Init(context.Context) error
	FetchBatch(context.Context, []string) ([]map[string]any, error)
	Close() error
}

type BatchError struct {
	Batch []string
	Err   error
}

func (e *BatchError) Error() string {
	return fmt.Sprintf("symbols %v: %v", e.Batch, e.Err)
}

func (e *BatchError) Unwrap() error {
	return e.Err
}

func RunPipeline(ctx context.Context, cfg PipelineConfig, fetcher BatchFetcher, log *zap.Logger, nc *nats.Conn) {

	log = log.Named(cfg.Type).With(
		zap.Int("batchSize", cfg.BatchSize),
		zap.Int("batchesPerTick", cfg.BatchesPerTick),
	)
	log.Info("pipeline started")

	if err := fetcher.Init(ctx); err != nil {
		log.Error("fetcher failed", zap.Error(err))
		return
	}
	defer func() { _ = fetcher.Close() }()

	tickEvery := cfg.Every
	if tickEvery <= 0 {
		tickEvery = 24 * time.Hour
	}
	ticker := time.NewTicker(tickEvery)
	defer ticker.Stop()

	tickCount := 0

	token := make(chan struct{}, 1)
	token <- struct{}{}

	job := func() {
		defer func() {
			if r := recover(); r != nil {
				log.Error("job panic", zap.Any("panic", r), zap.Stack("stack"))
			}
		}()

		timeOut := tickEvery - tickEvery/10
		jobCtx, cancel := context.WithTimeout(ctx, timeOut)
		defer cancel()

		merged, err := fetchPool(jobCtx, fetcher, cfg.Symbols, cfg.BatchSize, cfg.BatchesPerTick, cfg.MaxWorkers, cfg.MaxConcurrency)
		if err != nil {
			if jobCtx.Err() != nil || ctx.Err() != nil {
				log.Info("job canceled during shutdown", zap.Error(err))
				return
			}
			var be *BatchError
			if errors.As(err, &be) {
				log.Error("fetch error", zap.Error(err), zap.Strings("symbols", be.Batch))
			} else {
				log.Error("fetch error", zap.Error(err))
			}
			return
		}

		log.Debug("fetched merged results", zap.Int("items", len(merged.QuoteResponse.Result)))

		// Normalize and publish
		if err := usecase.NormalizeAndPublish(jobCtx, nc, merged, false, log); err != nil {
			log.Error("publish error", zap.Error(err))
		}

		write := cfg.WriteJson
		writeEvery := cfg.WriteJsonEveryTicks
		if writeEvery <= 0 {
			writeEvery = 1
		}
		if write {
			tickCount++
			if tickCount%writeEvery == 0 {
				fname := fmt.Sprintf("%s_%s.json", cfg.Type, time.Now().Format("20060102_150405"))
				fpath := filepath.Join(cfg.OutputDir, fname)
				if dir := filepath.Dir(fpath); dir != "" && dir != "." {
					_ = os.MkdirAll(dir, 0o755)
				}
				if err := os.WriteFile(fpath, writeJson(merged), 0o644); err != nil {
					log.Error("write file", zap.Error(err), zap.String("path", fpath))
					return
				}
				log.Debug("pipeline wrote file", zap.String("path", fpath), zap.Int("items", len(merged.QuoteResponse.Result)))
			}
		}

	}

	// Run job once immediately
	go func() {
		<-token
		defer func() {
			if r := recover(); r != nil {
				log.Error("job panic", zap.Any("panic", r), zap.Stack("stack"))
			}
			token <- struct{}{}
		}()
		job()
	}()

	// Main loop
	for {
		select {
		case <-ctx.Done():
			<-token
			if cause := context.Cause(ctx); cause != nil {
				log.Info("pipeline stopping", zap.Error(cause))
			} else {
				log.Info("pipeline stopping")
			}
			return
		case <-ticker.C:
			select {
			case <-ctx.Done():
				continue
			default:
			}
			select {
			case <-token:
				tickCount++
				go func() {
					defer func() {
						if r := recover(); r != nil {
							log.Error("job panic", zap.Any("panic", r), zap.Stack("stack"))
						}
						token <- struct{}{}
					}()
					job()
				}()
			default:
				log.Debug("skip tick: job busy", zap.String("type", cfg.Type))
			}
		}
	}
}

type QuoteStruct struct {
	QuoteResponse struct {
		Result []map[string]any `json:"result"`
	} `json:"quoteResponse"`
}

type job struct {
	idx   int
	batch []string
}

type result struct {
	idx   int
	batch []string
	items []map[string]any
	err   error
}

func fetchPool(
	ctx context.Context,
	fetcher BatchFetcher,
	symbols []string,
	batchSize int,
	batchesPerTick int,
	maxWorkers int,
	maxConcurrency int,
) (*QuoteStruct, error) {

	if batchSize <= 0 {
		batchSize = 10
	}
	if batchesPerTick <= 0 {
		batchesPerTick = 1
	}
	if maxWorkers <= 0 {
		maxWorkers = 2
	}
	if maxConcurrency <= 0 {
		maxConcurrency = 2
	}

	// Batch the symbols
	batches := splitBatches(symbols, batchSize)
	if len(batches) > batchesPerTick {
		batches = batches[:batchesPerTick]
	}

	// Declare channels and semaphore
	jobs := make(chan job)
	results := make(chan result, len(batches))
	sem := semaphore.NewWeighted(int64(maxConcurrency))
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := range jobs {
				if err := sem.Acquire(ctx, 1); err != nil {
					results <- result{idx: j.idx, batch: j.batch, err: fmt.Errorf("acquire semaphore: %w", err)}
					return
				}
				items, err := fetcher.FetchBatch(ctx, j.batch)
				sem.Release(1)
				results <- result{idx: j.idx, batch: j.batch, items: items, err: err}
			}
		}(i + 1)
	}

	// Distribute jobs
	go func() {
		for i, batch := range batches {
			select {
			case jobs <- job{idx: i, batch: batch}:
			case <-ctx.Done():
				close(jobs)
				return
			}
		}
		close(jobs)
	}()

	// Close results channel
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	merged := &QuoteStruct{}
	var firstErr error
	for r := range results {
		if r.err != nil && firstErr == nil {
			firstErr = &BatchError{Batch: append([]string(nil), r.batch...), Err: r.err}
		}
		if len(r.items) > 0 {
			merged.QuoteResponse.Result = append(merged.QuoteResponse.Result, r.items...)
		}
	}

	if firstErr != nil {
		return nil, firstErr
	}

	return merged, nil
}

func writeJson(v any) []byte { return util.WriteJson(v) }

func splitBatches(symbols []string, n int) [][]string {
	if n <= 0 || len(symbols) == 0 {
		return nil
	}
	out := make([][]string, 0, (len(symbols)+n-1)/n)
	for i := 0; i < len(symbols); i += n {
		j := i + n
		if j > len(symbols) {
			j = len(symbols)
		}
		out = append(out, symbols[i:j])
	}
	return out
}
