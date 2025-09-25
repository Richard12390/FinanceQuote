package domain

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type envelope struct {
	QuoteResponse struct {
		Result []map[string]any `json:"result"`
	} `json:"quoteResponse"`
}

func asString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case float64:
		return strconv.FormatFloat(t, 'f', -1, 64)
	case int64:
		return strconv.FormatInt(t, 10)
	case int:
		return strconv.FormatInt(int64(t), 10)
	case json.Number:
		return t.String()
	default:
		return ""
	}
}
func asInt64(v any) int64 {
	switch t := v.(type) {
	case string:
		x, _ := strconv.ParseInt(strings.TrimSpace(t), 10, 64)
		return x
	case float64:
		return int64(t)
	case int64:
		return t
	case int:
		return int64(t)
	case json.Number:
		x, _ := t.Int64()
		return x
	default:
		return 0
	}
}
func asFloat64(v any) float64 {
	switch t := v.(type) {
	case string:
		x, _ := strconv.ParseFloat(strings.TrimSpace(t), 64)
		return x
	case float64:
		return t
	case json.Number:
		x, _ := t.Float64()
		return x
	case int64:
		return float64(t)
	case int:
		return float64(t)
	default:
		return 0
	}
}
func has(m map[string]any, k string) bool { _, ok := m[k]; return ok }

func quoteBase(sym string) string {
	s := strings.ToUpper(sym)
	if strings.HasSuffix(s, "USDT") {
		return strings.TrimSuffix(s, "USDT")
	}
	if strings.HasSuffix(s, "USD") {
		return strings.TrimSuffix(s, "USD")
	}
	return s
}
func quoteSuffix(sym string) string {
	s := strings.ToUpper(sym)
	if strings.HasSuffix(s, "USDT") {
		return "USDT"
	}
	if strings.HasSuffix(s, "USD") {
		return "USD"
	}
	return ""
}

func fromBinance(rs []map[string]any) bool {
	for _, r := range rs {
		if !has(r, "symbol") {
			return false
		}
		if !(has(r, "lastPrice") || has(r, "price") || has(r, "weightedAvgPrice")) {
			return false
		}
	}
	return true
}
func fromYahoo(rs []map[string]any) bool {
	for _, r := range rs {
		if !(has(r, "quoteType") && has(r, "regularMarketPrice")) {
			return false
		}
	}
	return true
}
func fromTwse(rs []map[string]any) bool {
	cnt := 0
	for _, r := range rs {
		if has(r, "tlong") && (has(r, "z") || has(r, "@") || has(r, "ch")) {
			cnt++
		}
	}
	return cnt > 0 && cnt*2 >= len(rs)
}

func NormalizeEnvelope(in []byte, keepRaw bool) ([]QuoteNorm, error) {
	var env envelope
	if err := json.Unmarshal(in, &env); err != nil {
		return nil, err
	}
	rs := env.QuoteResponse.Result
	if len(rs) == 0 {
		return nil, nil
	}

	if fromBinance(rs) {
		out := make([]QuoteNorm, 0, len(rs))
		for _, r := range rs {
			out = append(out, normalizeBinance(r, keepRaw))
		}
		return out, nil
	}
	if fromYahoo(rs) {
		out := make([]QuoteNorm, 0, len(rs))
		for _, r := range rs {
			out = append(out, normalizeYahoo(r, keepRaw))
		}
		return out, nil
	}
	if fromTwse(rs) {
		out := make([]QuoteNorm, 0, len(rs))
		for _, r := range rs {
			out = append(out, normalizeTwse(r, keepRaw))
		}
		return out, nil
	}
	out := make([]QuoteNorm, 0, len(rs))
	for _, r := range rs {
		out = append(out, fallback(r, keepRaw))
	}
	return out, nil
}

func fallback(r map[string]any, keepRaw bool) QuoteNorm {
	ms := time.Now().UnixMilli()
	ts := time.UnixMilli(ms)
	out := QuoteNorm{TsUnixMs: ms, TsISO: ts.UTC().Format(time.RFC3339), TsLocal: ts.Local().Format(time.RFC3339)}
	if keepRaw {
		out.Raw = r
	}
	return out
}
