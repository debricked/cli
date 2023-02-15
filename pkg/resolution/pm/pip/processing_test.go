package pip

import (
	"testing"

	"github.com/debricked/cli/pkg/resolution/pm/writer"
	"github.com/stretchr/testify/assert"
)

func TestParsePipList(t *testing.T) {
	// TODO Fix test for parse Pip List
	job := NewJob("file", CmdFactory{}, writer.FileWriter{})
	pipData := "load pip list"
	job.parsePipList(pip)
	assert.Equal(t, "file", job.file)
	assert.Nil(t, job.err)
}

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
	metadata := "Load test-data"
	packages := []string{
		"Flask",
		"sentry-sdk",
		"sentry-sdk[flask]",
		"pandas",
	}

	nodes, edges, err := job.parseGraph(packages, metadata)
	assert.Equal(t, nodes, []string{"pandas 1.4.2", "numpy 1.21.5"})
	assert.Equal(t, edges, []string{"pandas python-dateutil", "pandas pytz", "pandas numpy"})
	assert.Nil(t, err)
	assert.Equal(t, "file", job.file)
	assert.Nil(t, job.err)
}

func TestParsePackageMetadata(t *testing.T) {
	// TODO Fix test for parse Package metadata
	job := NewJob("file", CmdFactory{}, writer.FileWriter{})
	assert.Equal(t, "file", job.file)
	assert.Nil(t, job.err)
}
