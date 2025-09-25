package symbolmap

import "strings"

func MapYahooTW(symbols []string) []string {
	res := make([]string, 0, len(symbols))
	for _, sym := range symbols {
		sym = strings.TrimSpace(sym)
		if sym == "" {
			continue
		}
		if strings.HasSuffix(sym, ".TW") {
			res = append(res, sym)
		} else {
			res = append(res, sym+".TW")
		}
	}
	return res
}

func MapYahooCrypto(symbols []string) []string {
	res := make([]string, 0, len(symbols))
	for _, sym := range symbols {
		sym = strings.ToUpper(strings.TrimSpace(sym))
		if strings.HasSuffix(sym, "USDT") {
			res = append(res, strings.TrimSuffix(sym, "USDT")+"-USD")
		} else if strings.HasSuffix(sym, "USD") {
			res = append(res, strings.TrimSuffix(sym, "USD")+"-USD")
		} else {
			res = append(res, sym+"-USD")
		}
	}
	return res
}
