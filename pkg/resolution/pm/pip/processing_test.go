package pip

import (
	"testing"

	"github.com/debricked/cli/pkg/resolution/pm/writer"
	"github.com/stretchr/testify/assert"
)

func TestParseRequirements(t *testing.T) {
	job := NewJob("testdata/requirements.txt", CmdFactory{}, writer.FileWriter{})
	packages, err := job.parseRequirements()
	assert.Nil(t, err)

	gt := []string{
		"Flask",
		"sentry-sdk",
		"sentry-sdk[flask]",
		"pandas",
		"tqdm",
		"cryptography",
	}
	assert.Equal(t, gt, packages)
	assert.Nil(t, job.err)
}

func TestParseGraph(t *testing.T) {
	// TODO Fix test for parse Graph
	job := NewJob("file", CmdFactory{}, writer.FileWriter{})
	assert.Equal(t, "file", job.file)
	assert.Nil(t, job.err)
}

func TestParsePipList(t *testing.T) {
	// TODO Fix test for parse Pip List
	job := NewJob("file", CmdFactory{}, writer.FileWriter{})
	assert.Equal(t, "file", job.file)
	assert.Nil(t, job.err)
}
