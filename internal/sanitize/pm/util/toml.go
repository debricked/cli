package util

import (
	"bytes"
	"encoding/json"
)

// QuoteTOML returns value as a TOML basic string. TOML and JSON strings use the same escapes.
func QuoteTOML(value string) string {
	var quoted bytes.Buffer
	encoder := json.NewEncoder(&quoted)
	encoder.SetEscapeHTML(false)
	_ = encoder.Encode(value) // a string always encodes

	return string(bytes.TrimRight(quoted.Bytes(), "\n"))
}
