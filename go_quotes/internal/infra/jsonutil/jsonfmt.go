package util

import (
	"bytes"
	"encoding/json"
)

func WriteJson(v any) []byte {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
	return buf.Bytes()
}
