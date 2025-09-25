package domain

import (
	"os"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	exchangeMapOnce      sync.Once
	yahooExchangeCodeMap map[string]string
)

func loadYahooExchangeMap() map[string]string {
	exchangeMapOnce.Do(func() {
		yahooExchangeCodeMap = make(map[string]string)

		candidates := []string{
			strings.TrimSpace(os.Getenv("YAHOO_EXCHANGE_MAP")),
			filepath.Join("configs", "exchange_map.yml"),
		}
		for _, path := range candidates {
			if path == "" {
				continue
			}
			data, err := os.ReadFile(path)
			if err != nil {
				continue
			}
			var fileMap map[string]string
			if err := yaml.Unmarshal(data, &fileMap); err != nil {
				continue
			}
			cleaned := make(map[string]string, len(fileMap))
			for k, v := range fileMap {
				key := strings.TrimSpace(k)
				val := strings.TrimSpace(v)
				if key == "" || val == "" {
					continue
				}
				cleaned[key] = val
			}
			if len(cleaned) > 0 {
				yahooExchangeCodeMap = cleaned
				break
			}
		}
	})
	return yahooExchangeCodeMap
}

func mapYahooExchange(name string) string {
	key := strings.TrimSpace(name)
	if key == "" {
		return ""
	}
	if code, ok := loadYahooExchangeMap()[key]; ok {
		return code
	}
	return key
}
