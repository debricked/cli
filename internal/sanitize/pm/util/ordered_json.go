package util

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"strings"
)

// OrderedObject is a JSON object that keeps its key order. encoding/json sorts map keys, and some
// backend parsers use the first or the last entry with a given name.
type OrderedObject struct {
	Keys   []string
	Values map[string]any
}

func NewOrderedObject() *OrderedObject {
	return &OrderedObject{Values: map[string]any{}}
}

func (o *OrderedObject) Set(key string, value any) {
	if _, exists := o.Values[key]; !exists {
		o.Keys = append(o.Keys, key)
	}
	o.Values[key] = value
}

func ParseJSONObject(content []byte) (*OrderedObject, error) {
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.UseNumber()
	value, err := decodeJSONValue(decoder)
	if err != nil {
		return nil, err
	}
	if _, err = decoder.Token(); !errors.Is(err, io.EOF) {
		return nil, errors.New("unexpected content after the JSON document")
	}
	object, isObject := value.(*OrderedObject)
	if !isObject {
		return nil, errors.New("the JSON document is not an object")
	}

	return object, nil
}

func decodeJSONValue(decoder *json.Decoder) (any, error) {
	token, err := decoder.Token()
	if err != nil {
		return nil, err
	}
	delimiter, isDelimiter := token.(json.Delim)
	if !isDelimiter {
		return token, nil
	}

	if delimiter == '[' {
		return decodeJSONArray(decoder)
	}

	return decodeJSONObject(decoder)
}

func decodeJSONArray(decoder *json.Decoder) (any, error) {
	array := []any{}
	for decoder.More() {
		value, err := decodeJSONValue(decoder)
		if err != nil {
			return nil, err
		}
		array = append(array, value)
	}
	_, err := decoder.Token()

	return array, err
}

func decodeJSONObject(decoder *json.Decoder) (any, error) {
	object := NewOrderedObject()
	for decoder.More() {
		keyToken, err := decoder.Token()
		if err != nil {
			return nil, err
		}
		key, isString := keyToken.(string)
		if !isString {
			return nil, errors.New("a JSON object key is not a string")
		}
		value, err := decodeJSONValue(decoder)
		if err != nil {
			return nil, err
		}
		object.Set(key, value)
	}
	_, err := decoder.Token()

	return object, err
}

func EncodeJSON(value any) ([]byte, error) {
	var buffer bytes.Buffer
	if err := writeJSONValue(&buffer, value, 0); err != nil {
		return nil, err
	}
	buffer.WriteByte('\n')

	return buffer.Bytes(), nil
}

func writeJSONValue(buffer *bytes.Buffer, value any, depth int) error {
	switch typed := value.(type) {
	case *OrderedObject:
		return writeJSONObject(buffer, typed, depth)
	case []any:
		return writeJSONArray(buffer, typed, depth)
	default:
		return writeJSONScalar(buffer, typed)
	}
}

func writeJSONObject(buffer *bytes.Buffer, object *OrderedObject, depth int) error {
	if len(object.Keys) == 0 {
		buffer.WriteString("{}")

		return nil
	}
	buffer.WriteByte('{')
	for index, key := range object.Keys {
		writeJSONLineStart(buffer, index, depth+1)
		if err := writeJSONScalar(buffer, key); err != nil {
			return err
		}
		buffer.WriteString(": ")
		if err := writeJSONValue(buffer, object.Values[key], depth+1); err != nil {
			return err
		}
	}
	writeJSONLineEnd(buffer, depth)
	buffer.WriteByte('}')

	return nil
}

func writeJSONArray(buffer *bytes.Buffer, array []any, depth int) error {
	if len(array) == 0 {
		buffer.WriteString("[]")

		return nil
	}
	buffer.WriteByte('[')
	for index, value := range array {
		writeJSONLineStart(buffer, index, depth+1)
		if err := writeJSONValue(buffer, value, depth+1); err != nil {
			return err
		}
	}
	writeJSONLineEnd(buffer, depth)
	buffer.WriteByte(']')

	return nil
}

func writeJSONLineStart(buffer *bytes.Buffer, index, depth int) {
	if index > 0 {
		buffer.WriteByte(',')
	}
	writeJSONLineEnd(buffer, depth)
}

func writeJSONLineEnd(buffer *bytes.Buffer, depth int) {
	buffer.WriteByte('\n')
	buffer.WriteString(strings.Repeat("    ", depth))
}

func writeJSONScalar(buffer *bytes.Buffer, value any) error {
	var scalar bytes.Buffer
	encoder := json.NewEncoder(&scalar)
	// Keeps version ranges such as ">=1.0" readable instead of writing ">=1.0".
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(value); err != nil {
		return err
	}
	buffer.Write(bytes.TrimRight(scalar.Bytes(), "\n"))

	return nil
}
