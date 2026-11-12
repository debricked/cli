package util

import "gopkg.in/yaml.v3"

// ParseYAMLMapping returns the top-level mapping of a YAML document with exactly one document.
func ParseYAMLMapping(content []byte) (*yaml.Node, error) {
	var document yaml.Node
	if err := yaml.Unmarshal(content, &document); err != nil || len(document.Content) != 1 {
		return nil, ErrInvalidFormat
	}
	root := document.Content[0]
	if root.Kind != yaml.MappingNode {
		return nil, ErrInvalidFormat
	}

	return root, nil
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
