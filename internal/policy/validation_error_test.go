package policy

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalPropertyPath(t *testing.T) {
	cases := []struct {
		name string
		json string
		path string
	}{
		{name: "context property", json: `{"message": "m", "context": {"property": "policies[0].name"}}`, path: "policies[0].name"},
		{name: "context pointer", json: `{"message": "m", "context": {"pointer": "/policies/0/name"}}`, path: "/policies/0/name"},
		{name: "context path segments", json: `{"message": "m", "context": {"path": ["policies", 0, "name"]}}`, path: "policies.0.name"},
		{name: "top level property path", json: `{"message": "m", "propertyPath": "policies[1].rules"}`, path: "policies[1].rules"},
		{name: "no path", json: `{"message": "m"}`, path: ""},
		{name: "empty path", json: `{"message": "m", "context": {"property": "  "}}`, path: ""},
		{name: "non string path", json: `{"message": "m", "context": {"property": 1}}`, path: ""},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var validationError ValidationError

			err := json.Unmarshal([]byte(c.json), &validationError)

			assert.NoError(t, err)
			assert.Equal(t, "m", validationError.Message)
			assert.Equal(t, c.path, validationError.PropertyPath)
		})
	}
}

func TestUnmarshalMalformedValidationError(t *testing.T) {
	var validationError ValidationError

	err := json.Unmarshal([]byte(`"not an object"`), &validationError)

	assert.Error(t, err)
}

func TestValidationErrorString(t *testing.T) {
	assert.Equal(t, "m", ValidationError{Message: "m"}.String())
	assert.Equal(t, "p: m", ValidationError{PropertyPath: "p", Message: "m"}.String())
	assert.Equal(t, "invalid value", ValidationError{}.String())
	assert.Equal(t, "p: invalid value", ValidationError{PropertyPath: "p"}.String())
}

func TestValidationErrorMarshal(t *testing.T) {
	validationError := ValidationError{
		PropertyPath: "policies[0].name",
		Message:      "The property name is required",
		Context:      map[string]interface{}{"constraint": "required"},
	}

	data, err := json.Marshal(validationError)

	assert.NoError(t, err)
	assert.JSONEq(
		t,
		`{"propertyPath": "policies[0].name", "message": "The property name is required", "context": {"constraint": "required"}}`,
		string(data),
	)
}

func TestJsonSyntaxError(t *testing.T) {
	assert.Nil(t, jsonSyntaxError([]byte(`{"version": "1.0"}`)))

	syntaxErr := jsonSyntaxError([]byte("{\n  \"version\": \"1.0\",,\n}"))
	assert.NotNil(t, syntaxErr)
	assert.Contains(t, syntaxErr.Message, "invalid JSON on line 2, column 21")

	typeErr := jsonSyntaxError([]byte(`{"version": 1.0e}`))
	assert.NotNil(t, typeErr)
	assert.Contains(t, typeErr.Message, "invalid JSON")
}

func TestPositionBeyondContent(t *testing.T) {
	line, column := position([]byte("ab"), 10)

	assert.Equal(t, 1, line)
	assert.Equal(t, 3, column)
}
