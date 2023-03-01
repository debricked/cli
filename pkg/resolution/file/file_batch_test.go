package file

import (
	"testing"

	"github.com/debricked/cli/pkg/resolution/pm/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewBatch(t *testing.T) {
	b := NewBatch(nil)
	assert.NotNil(t, b)

	b = NewBatch(testdata.PmMock{})
	assert.NotNil(t, b)
}

func TestFiles(t *testing.T) {
	b := NewBatch(testdata.PmMock{})

	files := b.Files()
	assert.Empty(t, files)

	b.Add("file-1")
	assert.Len(t, b.Files(), 1)

	b.Add("file-1")
	assert.Len(t, b.Files(), 1)

	b.Add("file-2")
	assert.Len(t, b.Files(), 2)
}

func TestAdd(t *testing.T) {
	b := NewBatch(testdata.PmMock{})

	filesMap := b.files
	assert.Empty(t, filesMap)

	b.Add("file-1")
	filesMap = b.files
	assert.Len(t, filesMap, 1)

	b.Add("file-1")
	filesMap = b.files
	assert.Len(t, filesMap, 1)

	b.Add("file-2")
	filesMap = b.files
	assert.Len(t, filesMap, 2)
}

func TestPm(t *testing.T) {
	b := NewBatch(nil)
	assert.Nil(t, b.Pm())

	pm := testdata.PmMock{}
	b = NewBatch(pm)
	assert.Equal(t, pm, b.Pm())
}
