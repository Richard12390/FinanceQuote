package yahoo

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type YahooFetcher struct {
	client *http.Client
	crumb  string
}

func (f *YahooFetcher) Init(ctx context.Context) error {
	c, crumb, err := getCookieAndCrumb()
	if err != nil {
		return err
	}
	f.client, f.crumb = c, crumb
	return nil
}

func (f *YahooFetcher) FetchBatch(ctx context.Context, batch []string) ([]map[string]any, error) {

	// u := fmt.Sprintf("https://query2.finance.yahoo.com/v7/finance/quote?symbols=%s&crumb=%s", strings.Join(batch, ","), y.crumb)

	base := &url.URL{Scheme: "https", Host: "query2.finance.yahoo.com", Path: "/v7/finance/quote"}
	q := base.Query()
	q.Set("symbols", strings.Join(batch, ","))
	q.Set("crumb", f.crumb)
	base.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, base.String(), nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json")

	resp, err := f.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Yahoo fetch: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Yahoo HTTP %d: %s", resp.StatusCode, string(b))
	}

	var env struct {
		QuoteResponse struct {
			Result []map[string]any `json:"result"`
			Error  any              `json:"error"`
		} `json:"quoteResponse"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, err
	}
	return env.QuoteResponse.Result, nil
}

func (f *YahooFetcher) Close() error { return nil }
