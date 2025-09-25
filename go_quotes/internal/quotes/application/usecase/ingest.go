package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go_quotes/internal/infra/ctxutil"
	"go_quotes/internal/quotes/domain"

	nats "github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func subjectFor(q domain.QuoteNorm, root string) string {
	if root == "" {
		root = "quotes"
	}
	src := strings.ToLower(strings.TrimSpace(q.Source))
	if src == "" {
		src = "unknown"
	}
	at := strings.ToLower(strings.TrimSpace(q.AssetType))
	if at == "" {
		at = "unknown"
	}
	return root + "." + src + "." + at
}

var EnablePersistAck = false

type persistReply struct {
	OK    *bool  `json:"ok"`
	Error string `json:"error,omitempty"`
}

func requestPersist(ctx context.Context, nc *nats.Conn, q domain.QuoteNorm) error {
	subject := "quotes.persist." + strings.ToLower(q.AssetType)
	body, _ := json.Marshal(q)
	c, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	msg, err := nc.RequestWithContext(c, subject, body)
	if err != nil {
		return fmt.Errorf("Persist request failed: %w", err)
	}
	var pr persistReply
	if err := json.Unmarshal(msg.Data, &pr); err != nil {
		return fmt.Errorf("Persist bad reply: %w", err)
	}
	if pr.OK == nil {
		return fmt.Errorf("Persist bad reply: missing 'ok'")
	}
	if !*pr.OK {
		return fmt.Errorf("Persist not ok: %s", pr.Error)
	}
	return nil
}

// NormalizeAndPublish： JSON → []QuoteNorm → NATS
func NormalizeAndPublish(ctx context.Context, nc *nats.Conn, merged any, keepRaw bool, log *zap.Logger) error {
	j, err := json.Marshal(merged)
	if err != nil {
		return fmt.Errorf("Marshal merged failed: %w", err)
	}

	norms, err := domain.NormalizeEnvelope(j, keepRaw)
	if err != nil {
		return fmt.Errorf("Normalize failed: %w", err)
	}
	if len(norms) == 0 {
		log.Debug("No quote to publish")
		return nil
	}

	for _, q := range norms {
		subject := subjectFor(q, "quotes")
		body, _ := json.Marshal(q)
		fmt.Println(subject, string(body))
		select {
		case <-ctx.Done():
			log.Debug("skip publish: shutting down")
			return nil
		default:
		}
		if err := nc.Publish(subject, body); err != nil {
			if ctxutil.IsCancel(err, ctx) {
				log.Info("NATS publish canceled during shutdown", zap.Error(err))
				return nil
			}
			return fmt.Errorf("NATS publish failed: %s: %w", string(body), err)
		}
		if EnablePersistAck {
			if err := requestPersist(ctx, nc, q); err != nil {
				return err
			}
		}
	}
	select {
	case <-ctx.Done():
		log.Debug("skip flush: shutting down")
		return nil
	default:
	}
	if err := nc.FlushTimeout(2 * time.Second); err != nil {
		if ctxutil.IsCancel(err, ctx) {
			log.Info("NATS flush canceled during shutdown", zap.Error(err))
			return nil
		}
		return fmt.Errorf("NATS flush timeout: %w", err)
	}
	// log.Debug("published", zap.String("subject", subjectFor(norms[0], "quotes")), zap.Int("count", len(norms)))
	return nil
}
