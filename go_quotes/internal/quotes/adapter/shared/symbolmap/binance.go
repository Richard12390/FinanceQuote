package symbolmap

import "strings"

func MapBinance(symbols []string) []string {
	res := make([]string, 0, len(symbols))
	for _, sym := range symbols {
		sym = strings.ToUpper(strings.TrimSpace(sym))
		if sym != "" {
			res = append(res, sym)
		}
	}
	return res
}
