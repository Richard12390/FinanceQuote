package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	nats "github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	cfg "go_quotes/internal/conf"
	cfgpath "go_quotes/internal/infra/configpath"
	"go_quotes/internal/quotes/adapter/outbound/fetcher/binance"
	"go_quotes/internal/quotes/adapter/outbound/fetcher/twse"
	"go_quotes/internal/quotes/adapter/outbound/fetcher/yahoo"
	sm "go_quotes/internal/quotes/adapter/shared/symbolmap"
	"go_quotes/internal/quotes/application/pipeline"

	"go_quotes/internal/infra/logging"
)

type YahooFetcher = yahoo.YahooFetcher
type TwseFetcher = twse.TwseFetcher
type BinanceRESTFetcher = binance.BinanceRestFetcher

func mapYahooTW(symbols []string) []string     { return sm.MapYahooTW(symbols) }
func mapYahooCrypto(symbols []string) []string { return sm.MapYahooCrypto(symbols) }
func mapTwse(symbols []string) []string        { return sm.MapTwse(symbols) }
func mapBinance(symbols []string) []string     { return sm.MapBinance(symbols) }

func main() {
	configFlag := flag.String("config", "", "config file or directory(optional)")
	flag.Parse()

	found := cfgpath.Resolve(*configFlag)
	fmt.Println(found)

	if found == "" {
		panic(fmt.Sprintf("config file not found in %s", *configFlag))
	}

	// Load config
	config, err := cfg.Load(found)
	if err != nil {
		panic(fmt.Errorf("load config: %w", err))
	}

	// Set up NATS
	natsUrl := config.NATSUrl
	if env := os.Getenv("NATS_URL"); env != "" {
		natsUrl = env
	}
	if natsUrl == "" {
		natsUrl = nats.DefaultURL
	}

	// NATS connection
	nc, err := nats.Connect(
		natsUrl,
		nats.Name("quote-publisher"),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		panic(err)
	}
	defer nc.Drain()

	// Context
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Set up stock loggers
	stockLogger, stockRotator, err := logging.SetupLogger("logs/stock.log", zapcore.InfoLevel, 50, 7, 1)
	if err != nil {
		panic(err)
	}
	defer stockLogger.Sync()
	logging.StartDailyRotate(ctx, stockRotator, time.Local)

	// Set up ETF loggers
	etfLogger, etfRotator, err := logging.SetupLogger("logs/etf.log", zapcore.InfoLevel, 50, 7, 1)
	if err != nil {
		panic(err)
	}
	defer etfLogger.Sync()
	logging.StartDailyRotate(ctx, etfRotator, time.Local)

	// Set up crypto loggers
	cryptoLogger, cryptoRotator, err := logging.SetupLogger("logs/crypto.log", zapcore.InfoLevel, 50, 7, 1)
	if err != nil {
		panic(err)
	}
	defer cryptoLogger.Sync()
	logging.StartDailyRotate(ctx, cryptoRotator, time.Local)

	// Name loggers
	stockLogger = stockLogger.Named("stock")
	etfLogger = etfLogger.Named("etf")
	cryptoLogger = cryptoLogger.Named("crypto")

	// Set up fetchers
	yfStock, yEtf, yfCrypto := &YahooFetcher{}, &YahooFetcher{}, &YahooFetcher{}
	tfStock, tfEtf := &TwseFetcher{}, &TwseFetcher{}
	bfCrypto := &BinanceRESTFetcher{}

	// Run pipelines
	var wg sync.WaitGroup
	for _, p := range config.Pipelines {

		// Map symbols
		var syms []string
		switch p.Source + "/" + p.AssetType {
		case "yahoo/stock":
			syms = mapYahooTW(p.Symbols)
		case "yahoo/etf":
			syms = mapYahooTW(p.Symbols)
		case "yahoo/crypto":
			syms = mapYahooCrypto(p.Symbols)
		case "twse/stock":
			syms = mapTwse(p.Symbols)
		case "twse/etf":
			syms = mapTwse(p.Symbols)
		case "binance/crypto":
			syms = mapBinance(p.Symbols)
		default:
			continue
		}

		// Bind fetchers and loggers
		var fetcher pipeline.BatchFetcher
		var logger = stockLogger
		switch p.Source {
		case "yahoo":
			switch p.AssetType {
			case "stock":
				fetcher = yfStock
				logger = stockLogger
			case "etf":
				fetcher = yEtf
				logger = etfLogger
			case "crypto":
				fetcher = yfCrypto
				logger = cryptoLogger
			default:
				continue
			}
		case "twse":
			switch p.AssetType {
			case "stock":
				fetcher = tfStock
				logger = stockLogger
			case "etf":
				fetcher = tfEtf
				logger = etfLogger
			default:
				continue
			}
		case "binance":
			switch p.AssetType {
			case "crypto":
				fetcher = bfCrypto
				logger = cryptoLogger
			default:
				continue
			}
		default:
			continue
		}

		// Pipeline parameters
		write := false
		if p.WriteJson != nil {
			write = *p.WriteJson
		}
		writeEvery := p.WriteJsonEveryTicks
		if write && writeEvery <= 0 {
			writeEvery = 1
		}
		every, err := time.ParseDuration(p.Every)
		if err != nil {
			every = time.Second
		}
		pc := pipeline.PipelineConfig{
			Type:                p.Name,
			Symbols:             syms,
			Every:               every,
			BatchSize:           p.BatchSize,
			BatchesPerTick:      p.BatchesPerTick,
			MaxWorkers:          p.MaxWorkers,
			MaxConcurrency:      p.MaxConcurrency,
			OutputDir:           p.OutputDir,
			WriteJson:           write,
			WriteJsonEveryTicks: writeEvery,
		}
		wg.Add(1)
		go func(pc pipeline.PipelineConfig, f pipeline.BatchFetcher, log *zap.Logger) {
			defer wg.Done()
			pipeline.RunPipeline(ctx, pc, f, log, nc)
		}(pc, fetcher, logger)
	}

	<-ctx.Done()
	wg.Wait()
	stockLogger.Info("Shutting down")
	etfLogger.Info("Shutting down")
	cryptoLogger.Info("Shutting down")
}
