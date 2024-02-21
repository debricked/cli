package fingerprint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateWFP(t *testing.T) {

	minFileSize := 1
	w := NewWinnowing(&minFileSize)

	snippets, err := w.GenerateWFP("testdata/snippet/main.py")
	assert.NoError(t, err)
	assert.Equal(t, 8, len(*snippets))
	assert.Equal(t, "5e6ddca9", (*snippets)[0].Hash)
	assert.Equal(t, 14, (*snippets)[0].Line)
	assert.Equal(t, "def test():\n    print(\"Hello, World!\")\n\n\ndef test2():\n    print(\"Hello, World!2\")\n\n\ndef test3():\n    print(\"Hello, World!3\")\n\n\ndef test4():\n    print(\"Hello, Worl", (*snippets)[0].Content)

}
