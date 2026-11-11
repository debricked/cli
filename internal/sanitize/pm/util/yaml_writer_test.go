package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeYAML(t *testing.T) {
	inner := NewOrderedObject()
	inner.Set("version", "1.10")
	inner.Set("dependencies", NewOrderedObject())
	object := NewOrderedObject()
	object.Set("lockfileVersion", RawYAMLScalar("6.0"))
	object.Set("/express@4.18.2", inner)
	object.Set("name", "plain-name")

	assert.Equal(t, `lockfileVersion: 6.0
/express@4.18.2:
  version: '1.10'
  dependencies: {}
name: plain-name
`, string(EncodeYAML(object)))
}

func TestYAMLScalar(t *testing.T) {
	cases := map[string]string{
		"plain name":           "express",
		"version with 3 parts": "4.18.2",
		"short version":        "'1.10'",
		"integer":              "'42'",
		"boolean word":         "'true'",
		"yes in capitals":      "'YES'",
		"scoped name":          "'@acme/logger'",
		"single quote":         "'it''s'",
		"newline":              `"a\nb"`,
	}
	values := map[string]string{
		"plain name":           "express",
		"version with 3 parts": "4.18.2",
		"short version":        "1.10",
		"integer":              "42",
		"boolean word":         "true",
		"yes in capitals":      "YES",
		"scoped name":          "@acme/logger",
		"single quote":         "it's",
		"newline":              "a\nb",
	}
	for name, expected := range cases {
		assert.Equal(t, expected, yamlScalar(values[name]), name)
	}
}
