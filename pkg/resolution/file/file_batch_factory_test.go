package file

import (
	"github.com/debricked/cli/pkg/resolution/pm"
	"github.com/debricked/cli/pkg/resolution/pm/testdata"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewBatchFactory(t *testing.T) {
	bf := NewBatchFactory()
	assert.NotNil(t, bf)

	pms := bf.pms
	assert.Equal(t, pm.Pms(), pms)
}

func TestMakeNoPms(t *testing.T) {
	bf := BatchFactory{}
	batches := bf.Make([]string{"go.mod"})
	assert.Empty(t, batches)
}

func TestMakeNoFiles(t *testing.T) {
	bf := NewBatchFactory()
	batches := bf.Make([]string{})
	assert.Empty(t, batches)
}

func TestMakeNoManifests(t *testing.T) {
	bf := BatchFactory{pms: []pm.IPm{testdata.PmMock{}}}
	batches := bf.Make([]string{"go.mod"})
	assert.Empty(t, batches)
}

func TestMakeOneFile(t *testing.T) {
	bf := BatchFactory{pms: []pm.IPm{
		testdata.PmMock{
			N:  "go",
			Ms: []string{"go.mod"},
		},
	}}
	batches := bf.Make([]string{"test/go.mod"})
	assert.Len(t, batches, 1)
	batch := batches[0]
	assert.Len(t, batch.Files(), 1)
	file := batch.Files()[0]
	assert.Equal(t, "test/go.mod", file)
}

func TestMakeMultipleFiles(t *testing.T) {
	bf := BatchFactory{pms: []pm.IPm{
		testdata.PmMock{
			N:  "go",
			Ms: []string{"go.mod"},
		},
	}}
	batches := bf.Make([]string{"go.mod", "test/go.mod"})
	assert.Len(t, batches, 1)
	batch := batches[0]
	assert.Len(t, batch.Files(), 2)
}

func TestMakeMultipleBatches(t *testing.T) {
	bf := BatchFactory{pms: []pm.IPm{
		testdata.PmMock{
			N:  "go",
			Ms: []string{"go.mod"},
		},
		testdata.PmMock{
			N:  "mvn",
			Ms: []string{"pom.xml"},
		},
	}}
	batches := bf.Make([]string{"go.mod", "test/pom.xml"})
	assert.Len(t, batches, 2)
	for _, batch := range batches {
		assert.Len(t, batch.Files(), 1)
	}
}

func TestMakeMultipleBatchesMultipleFiles(t *testing.T) {
	bf := BatchFactory{pms: []pm.IPm{
		testdata.PmMock{
			N:  "go",
			Ms: []string{"go.mod"},
		},
		testdata.PmMock{
			N:  "mvn",
			Ms: []string{"pom.xml"},
		},
	}}
	batches := bf.Make([]string{"go.mod", "test/pom.xml", "test/sub/go.mod"})
	assert.Len(t, batches, 2)
	for _, batch := range batches {
		if len(batch.Files()) == 1 {
			assert.Contains(t, batch.Files(), "test/pom.xml")
		} else if len(batch.Files()) == 2 {
			assert.Contains(t, batch.Files(), "go.mod")
			assert.Contains(t, batch.Files(), "test/sub/go.mod")
		} else {
			t.Error("failed to assert number of files in the batch")
		}
	}
}
