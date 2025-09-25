package symbolmap

import (
	"fmt"
	"strings"
)

func MapTwse(symbols []string) []string {
	res := make([]string, 0, len(symbols))
	for _, sym := range symbols {
		sym = strings.TrimSpace(sym)
		if sym != "" {
			res = append(res, fmt.Sprintf("tse_%s.tw", sym))
		}
	}
	return res
}
