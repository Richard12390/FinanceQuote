package domain

import (
	"strings"
	"time"
)

func normalizeBinance(r map[string]any, keepRaw bool) QuoteNorm {
	symbolOrigin := strings.ToUpper(asString(r["symbol"]))
	base := quoteBase(symbolOrigin)
	quote := quoteSuffix(symbolOrigin)
	symbol := base
	if quote != "" {
		symbol += "-" + quote
	}
	ms := asInt64(r["closeTime"])
	if ms == 0 {
		ms = time.Now().UnixMilli()
	}

	tsms := time.UnixMilli(ms)
	tsIso := tsms.UTC().Format(time.RFC3339)
	tsLocal := tsms.Local().Format(time.RFC3339)

	out := QuoteNorm{
		Source:       "binance",
		AssetType:    "crypto",
		Symbol:       symbol,
		SymbolOrigin: symbolOrigin,
		Exchange:     "BINANCE",
		Currency:     quote,
		Price:        asFloat64(r["lastPrice"]),
		Open:         asFloat64(r["openPrice"]),
		High:         asFloat64(r["highPrice"]),
		Low:          asFloat64(r["lowPrice"]),
		PrevClose:    asFloat64(r["prevClosePrice"]),
		Bid:          asFloat64(r["bidPrice"]),
		BidSize:      asFloat64(r["askPrice"]),
		Ask:          asFloat64(r["bidQty"]),
		AskSize:      asFloat64(r["askQty"]),
		Volume:       asFloat64(r["volume"]),
		VolumeKind:   "24hr",
		TsUnixMs:     ms,
		TsISO:        tsIso,
		TsLocal:      tsLocal,
	}
	if keepRaw {
		out.Raw = r
	}
	return out
}
