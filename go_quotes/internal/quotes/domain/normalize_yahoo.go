package domain

import (
	"strings"
	"time"
)

func normalizeYahoo(r map[string]any, keepRaw bool) QuoteNorm {
	qt := strings.ToUpper(asString(r["quoteType"]))
	asset := "stock"
	switch qt {
	case "CRYPTOCURRENCY":
		asset = "crypto"
	case "ETF":
		asset = "etf"
	}

	vol := asFloat64(r["regularMarketVolume"])
	volKind := "regular"
	if asset == "crypto" {
		if v24 := asFloat64(r["volume24Hr"]); v24 > 0 {
			vol = v24
			volKind = "24hr"
		}
	}

	ms := (asInt64(r["regularMarketTime"])) * 1000
	if ms == 0 {
		ms = time.Now().UnixMilli()
	}
	tsms := time.UnixMilli(ms)
	tsIso := tsms.UTC().Format(time.RFC3339)
	tsLocal := tsms.Local().Format(time.RFC3339)

	out := QuoteNorm{
		Source:       "yahoo",
		AssetType:    asset,
		Symbol:       asString(r["symbol"]),
		SymbolOrigin: asString(r["symbol"]),
		Exchange:     mapYahooExchange(asString(r["fullExchangeName"])),
		Currency:     asString(r["currency"]),
		Price:        asFloat64(r["regularMarketPrice"]),
		Open:         asFloat64(r["regularMarketOpen"]),
		High:         asFloat64(r["regularMarketDayHigh"]),
		Low:          asFloat64(r["regularMarketDayLow"]),
		PrevClose:    asFloat64(r["regularMarketPreviousClose"]),
		Bid:          asFloat64(r["bid"]),
		BidSize:      asFloat64(r["ask"]),
		Ask:          asFloat64(r["bidSize"]),
		AskSize:      asFloat64(r["askSize"]),
		Volume:       vol,
		VolumeKind:   volKind,
		MarketState:  asString(r["marketState"]),
		TsUnixMs:     ms,
		TsISO:        tsIso,
		TsLocal:      tsLocal,
	}
	if keepRaw {
		out.Raw = r
	}
	return out
}
