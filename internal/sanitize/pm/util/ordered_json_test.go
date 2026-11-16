package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAndEncodeJSONKeepsKeyOrderAndValues(t *testing.T) {
	content := `{"z": 1.10, "a": {"b": [true, null, "x", []], "c": {}}, "html": ">=1.0 <2 &"}`

	object, err := ParseJSONObject([]byte(content))
	require.NoError(t, err)
	encoded, err := EncodeJSON(object)

	assert.NoError(t, err)
	assert.Equal(t, `{
    "z": 1.10,
    "a": {
        "b": [
            true,
            null,
            "x",
            []
        ],
        "c": {}
    },
    "html": ">=1.0 <2 &"
}
`, string(encoded))
}

func TestOrderedObjectSetKeepsTheFirstPositionAndTheLastValue(t *testing.T) {
	object := NewOrderedObject()
	object.Set("a", 1)
	object.Set("b", 2)
	object.Set("a", 3)

	assert.Equal(t, []string{"a", "b"}, object.Keys)
	assert.Equal(t, 3, object.Values["a"])
}

func TestParseJSONObjectInvalid(t *testing.T) {
	for _, content := range []string{
		``,
		`[]`,
		`"text"`,
		`{} {}`,
		`{"a": 1`,
		`{"a": [1, }`,
		`{"a": {"b": }}`,
		`{"a" 1}`,
	} {
		_, err := ParseJSONObject([]byte(content))
		assert.Error(t, err, content)
	}
}

func TestEncodeJSONFailsOnValuesJSONCannotHold(t *testing.T) {
	inObject := NewOrderedObject()
	inObject.Set("a", func() {})

	for name, value := range map[string]any{
		"scalar":    make(chan int),
		"in object": inObject,
		"in array":  []any{"ok", func() {}},
	} {
		_, err := EncodeJSON(value)
		assert.Error(t, err, name)
	}
}
