package util

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

// RawYAMLScalar is written as it is. It keeps a number in the source file a number.
type RawYAMLScalar string

var (
	yamlPlainRegex   = regexp.MustCompile(`^[/A-Za-z_.][\w@/.+~-]*$`)
	yamlVersionRegex = regexp.MustCompile(`^\d[\w+~-]*(?:\.[\w+~-]+){2,}$`)
	yamlReserved     = map[string]bool{
		"y": true, "n": true, "yes": true, "no": true, "on": true, "off": true,
		"true": true, "false": true, "null": true, ".inf": true, ".nan": true,
	}
)

// EncodeYAML writes nested objects of strings. yaml.v3 can't be used for this: it writes keys over
// 128 characters in the `? key` form, pnpm keys with peer suffixes often reach that length, and not
// all YAML parsers support the form.
func EncodeYAML(object *OrderedObject) []byte {
	var buffer bytes.Buffer
	writeYAMLObject(&buffer, object, 0)

	return buffer.Bytes()
}

func writeYAMLObject(buffer *bytes.Buffer, object *OrderedObject, depth int) {
	indent := strings.Repeat("  ", depth)
	for _, key := range object.Keys {
		buffer.WriteString(indent + yamlScalar(key) + ":")
		switch value := object.Values[key].(type) {
		case *OrderedObject:
			if len(value.Keys) == 0 {
				buffer.WriteString(" {}\n")

				continue
			}
			buffer.WriteByte('\n')
			writeYAMLObject(buffer, value, depth+1)
		case RawYAMLScalar:
			buffer.WriteString(" " + string(value) + "\n")
		case string:
			buffer.WriteString(" " + yamlScalar(value) + "\n")
		}
	}
}

// yamlScalar leaves a value unquoted only when it can't be read as anything but a string, because
// the backend parser fails on non-string versions and an unquoted 1.10 would lose its last zero.
func yamlScalar(value string) string {
	if strings.ContainsFunc(value, func(r rune) bool { return r < ' ' }) {
		encoded, _ := json.Marshal(value)

		return string(encoded)
	}
	if isPlainYAMLString(value) {
		return value
	}

	return "'" + strings.ReplaceAll(value, "'", "''") + "'"
}

func isPlainYAMLString(value string) bool {
	if yamlVersionRegex.MatchString(value) {
		return true
	}
	if !yamlPlainRegex.MatchString(value) || yamlReserved[strings.ToLower(value)] {
		return false
	}
	_, err := strconv.ParseFloat(value, 64)

	return err != nil
}
