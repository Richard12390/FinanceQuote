package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type BinanceRestFetcher struct{}

func (b *BinanceRestFetcher) Init(ctx context.Context) error { return nil }

func (b *BinanceRestFetcher) FetchBatch(ctx context.Context, batch []string) ([]map[string]any, error) {
	arr, _ := json.Marshal(batch)
	base, err := url.Parse("https://api.binance.com/api/v3/ticker/24hr") // 用 24hr 端點才有 volume
	if err != nil {
		return nil, err
	}
	q := base.Query()
	q.Set("symbols", string(arr))
	base.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, base.String(), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Binance fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= http.StatusBadRequest {
		bts, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Binance HTTP %d: %s", resp.StatusCode, string(bts))
	}
	var result []map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (b *BinanceRestFetcher) Close() error { return nil }
