package domain

import (
	"strings"
	"time"
)

func mapTwseExchange(code string) string {
	switch strings.ToUpper(strings.TrimSpace(code)) {
	case "TSE":
		return "TWSE"
	case "OTC":
		return "TWO"
	default:
		return strings.ToUpper(strings.TrimSpace(code))
	}
}

func normalizeTwse(r map[string]any, keepRaw bool) QuoteNorm {
	symbolOrigin := asString(r["@"])
	if symbolOrigin == "" {
		symbolOrigin = asString(r["ch"])
	}
	symbol := strings.ToUpper(strings.ReplaceAll(symbolOrigin, ".tw", ".TW"))

	asset := "stock"
	if asString(r["nu"]) != "" {
		asset = "etf"
	}

	exchange := mapTwseExchange(asString(r["ex"]))

	ms := (asInt64(r["tlong"])) * 1000
	if ms == 0 {
		ms = time.Now().UnixMilli()
	}
	tsms := time.UnixMilli(ms)
	tsIso := tsms.UTC().Format(time.RFC3339)
	tsLocal := tsms.Local().Format(time.RFC3339)

	out := QuoteNorm{
		Source:       "twse",
		AssetType:    asset,
		Symbol:       symbol,
		SymbolOrigin: symbolOrigin,
		Exchange:     exchange,
		Currency:     "TWD",
		Price:        asFloat64(r["z"]),
		Open:         asFloat64(r["o"]),
		High:         asFloat64(r["h"]),
		Low:          asFloat64(r["l"]),
		PrevClose:    asFloat64(r["v"]),
		Bid:          asFloat64(r["b"]),
		BidSize:      asFloat64(r["g"]),
		Ask:          asFloat64(r["a"]),
		AskSize:      asFloat64(r["f"]),
		Volume:       asFloat64(r["v"]),
		VolumeKind:   "regular",
		TsUnixMs:     ms,
		TsISO:        tsIso,
		TsLocal:      tsLocal,
	}
	if keepRaw {
		out.Raw = r
	}
	return out
}
