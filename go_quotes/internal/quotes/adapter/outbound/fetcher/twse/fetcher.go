package twse

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type TwseFetcher struct{}

func (f *TwseFetcher) Init(ctx context.Context) error { return nil }

func (t *TwseFetcher) FetchBatch(ctx context.Context, batch []string) ([]map[string]any, error) {
	u := "https://mis.twse.com.tw/stock/api/getStockInfo.jsp"
	q := url.Values{}
	q.Set("json", "1")
	q.Set("delay", "0")
	q.Set("ex_ch", strings.Join(batch, "|"))
	q.Set("_", fmt.Sprintf("%d", time.Now().UnixMilli()))
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u+"?"+q.Encode(), nil)

	req.Header.Set("Referer", "https://mis.twse.com.tw/stock/index.jsp")
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("TWSE fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		bts, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("TWSEFetcher HTTP %d: %s", resp.StatusCode, string(bts))
	}
	var result struct {
		MsgArray []map[string]any `json:"msgArray"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return result.MsgArray, nil
}

func (t *TwseFetcher) Close() error { return nil }
