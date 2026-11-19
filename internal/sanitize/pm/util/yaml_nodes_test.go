package util

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseYAMLMappingRejectsSeveralDocuments(t *testing.T) {
	_, err := ParseYAMLMapping([]byte("a: 1\n---\nb: 2\n"))

	assert.ErrorIs(t, err, ErrInvalidFormat)
}

func TestParseYAMLMappings(t *testing.T) {
	documents, err := ParseYAMLMappings([]byte("---\na: 1\n---\nb: 2\n"))

	assert.NoError(t, err)
	assert.Len(t, documents, 2)
	assert.Equal(t, "b", documents[1].Content[0].Value)

	for _, invalid := range []string{"", "a: [", "a: 1\n---\n- b\n"} {
		_, err = ParseYAMLMappings([]byte(invalid))
		assert.ErrorIs(t, err, ErrInvalidFormat, invalid)
	}
}

func TestReasonReportsOnlyFixedMessages(t *testing.T) {
	assert.Equal(t, ErrNoLockDocument.Error(), Reason(ErrNoLockDocument))
	assert.Equal(t, ErrUnsupportedVersion.Error(), Reason(ErrUnsupportedVersion))
	assert.Equal(t, ErrInvalidFormat.Error(), Reason(errors.New("line 3: s3cret")))
}
