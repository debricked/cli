package util

import (
	"bytes"
	"errors"
	"io"

	"gopkg.in/yaml.v3"
)

// ParseYAMLMapping returns the top-level mapping of a file with exactly one YAML document. A file
// with more documents fails, rather than having all but the first silently dropped.
func ParseYAMLMapping(content []byte) (*yaml.Node, error) {
	documents, err := ParseYAMLMappings(content)
	if err != nil || len(documents) != 1 {
		return nil, ErrInvalidFormat
	}

	return documents[0], nil
}

// ParseYAMLMappings returns the top-level mapping of each document in a file, in order.
func ParseYAMLMappings(content []byte) ([]*yaml.Node, error) {
	decoder := yaml.NewDecoder(bytes.NewReader(content))
	var mappings []*yaml.Node
	for {
		var document yaml.Node
		err := decoder.Decode(&document)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil || len(document.Content) != 1 || document.Content[0].Kind != yaml.MappingNode {
			return nil, ErrInvalidFormat
		}
		mappings = append(mappings, document.Content[0])
	}
	if len(mappings) == 0 {
		return nil, ErrInvalidFormat
	}

	return mappings, nil
}

// ForEachYAMLPair visits the pairs of a mapping that have a scalar key, in order.
func ForEachYAMLPair(mapping *yaml.Node, visit func(key string, value *yaml.Node)) {
	if mapping.Kind != yaml.MappingNode {
		return
	}
	for index := 0; index+1 < len(mapping.Content); index += 2 {
		if mapping.Content[index].Kind == yaml.ScalarNode {
			visit(mapping.Content[index].Value, mapping.Content[index+1])
		}
	}
}

// YAMLValue returns the value of the last pair with the given key, or nil.
func YAMLValue(mapping *yaml.Node, key string) *yaml.Node {
	var found *yaml.Node
	ForEachYAMLPair(mapping, func(candidate string, value *yaml.Node) {
		if candidate == key {
			found = value
		}
	})

	return found
}

func IsYAMLNull(node *yaml.Node) bool {
	return node.Kind == yaml.ScalarNode && node.Tag == "!!null"
}

// YAMLString returns the value of a non-null scalar.
func YAMLString(node *yaml.Node) (string, bool) {
	if node == nil || node.Kind != yaml.ScalarNode || IsYAMLNull(node) {
		return "", false
	}

	return node.Value, true
}
