package policy

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// propertyPathKeys are the keys the backend may use to point out the offending
// property, either on the error itself or inside its context.
var propertyPathKeys = []string{
	"propertyPath",
	"property_path",
	"property",
	"dataPointer",
	"pointer",
	"path",
	"field",
}

// ValidationError is one schema violation reported for a policy file
type ValidationError struct {
	PropertyPath string                 `json:"propertyPath,omitempty"`
	Message      string                 `json:"message"`
	Context      map[string]interface{} `json:"context,omitempty"`
}

func (e *ValidationError) UnmarshalJSON(data []byte) error {
	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	if message, ok := fields["message"].(string); ok {
		e.Message = message
	}
	if context, ok := fields["context"].(map[string]interface{}); ok {
		e.Context = context
	}
	e.PropertyPath = propertyPath(fields, e.Context)

	return nil
}

// String renders the error on one line, prefixed with the property path when
// the backend supplied one.
func (e ValidationError) String() string {
	message := e.Message
	if len(message) == 0 {
		message = "invalid value"
	}
	if len(e.PropertyPath) == 0 {
		return message
	}

	return fmt.Sprintf("%s: %s", e.PropertyPath, message)
}

func propertyPath(fields ...map[string]interface{}) string {
	for _, field := range fields {
		for _, key := range propertyPathKeys {
			if path := pathValue(field[key]); len(path) > 0 {
				return path
			}
		}
	}

	return ""
}

// pathValue normalizes a property path, which is either a string or a list of
// path segments depending on the validator used by the backend.
func pathValue(value interface{}) string {
	switch path := value.(type) {
	case string:
		return strings.TrimSpace(path)
	case []interface{}:
		var segments []string
		for _, segment := range path {
			segments = append(segments, fmt.Sprintf("%v", segment))
		}

		return strings.Join(segments, ".")
	}

	return ""
}

// jsonSyntaxError turns malformed JSON into a validation error pointing at the
// line and column that broke the parse.
func jsonSyntaxError(content []byte) *ValidationError {
	var document json.RawMessage
	err := json.Unmarshal(content, &document)
	if err == nil {
		return nil
	}

	message := fmt.Sprintf("invalid JSON. %s", err)
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		line, column := position(content, syntaxErr.Offset)
		message = fmt.Sprintf("invalid JSON on line %d, column %d. %s", line, column, syntaxErr)
	}

	return &ValidationError{Message: message}
}

func position(content []byte, offset int64) (int, int) {
	if offset > int64(len(content)) {
		offset = int64(len(content))
	}
	line := 1
	column := 1
	for _, char := range content[:offset] {
		if char == '\n' {
			line++
			column = 1
		} else {
			column++
		}
	}

	return line, column
}
