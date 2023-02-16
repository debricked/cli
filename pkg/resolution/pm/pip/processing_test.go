package pip

import (
	"os"
	"testing"

	"github.com/debricked/cli/pkg/resolution/pm/writer"
	"github.com/stretchr/testify/assert"
)

func TestParsePipList(t *testing.T) {
	job := NewJob("file", false, CmdFactory{}, writer.FileWriter{})
	file, err := os.ReadFile("testdata/list.txt")
	assert.Nil(t, err)
	pipData := string(file)
	packages, err := job.parsePipList(pipData)
	assert.Nil(t, err)
	gt := []string{"aiohttp", "cryptography", "numpy", "Flask", "open-source-health", "pandas", "tqdm"}
	assert.Equal(t, gt, packages)
	assert.Nil(t, job.err)
}

func TestParseRequirements(t *testing.T) {
	job := NewJob("testdata/requirements.txt", false, CmdFactory{}, writer.FileWriter{})
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
	job := NewJob("file", false, CmdFactory{}, writer.FileWriter{})
	output, err := os.ReadFile("testdata/show.txt")
	assert.Nil(t, err)
	metadata := string(output)
	packages := []string{
		"Flask",
		"sentry-sdk",
		"sentry-sdk[flask]",
		"pandas",
	}

	nodes, edges, missed, err := job.parseGraph(packages, metadata)
	gtNodes := []string{"Flask 2.1.2", "pandas 1.4.2", "numpy 1.21.5"}
	gtEdges := []string{
		"Flask click",
		"Flask importlib-metadata",
		"Flask itsdangerous",
		"Flask Jinja2",
		"Flask Werkzeug",
		"pandas python-dateutil",
		"pandas pytz",
		"pandas numpy",
	}
	// More missed than usual since show-file is very empty
	gtMissed := []string{"sentry-sdk", "sentry-sdk[flask]", "click", "importlib-metadata", "itsdangerous", "jinja2", "werkzeug", "python-dateutil", "pytz"}
	assert.Equal(t, gtNodes, nodes)
	assert.Equal(t, gtEdges, edges)
	assert.Equal(t, gtMissed, missed)
	assert.Nil(t, err)
	assert.Nil(t, job.err)
}

func TestParsePackageMetadata(t *testing.T) {
	job := NewJob("file", false, CmdFactory{}, writer.FileWriter{})
	output, err := os.ReadFile("testdata/show.txt")
	assert.Nil(t, err)
	showMetadata := string(output)

	packageMetadata, err := job.parsePackageMetadata(showMetadata)
	assert.Nil(t, err)

	gt := map[string]PackageMetadata{
		"flask": {
			Name:         "Flask",
			Version:      "2.1.2",
			Dependencies: []string{"click", "importlib-metadata", "itsdangerous", "Jinja2", "Werkzeug"},
		},
		"numpy": {
			Name:         "numpy",
			Version:      "1.21.5",
			Dependencies: []string{},
		},
		"pandas": {
			Name:         "pandas",
			Version:      "1.4.2",
			Dependencies: []string{"python-dateutil", "pytz", "numpy"},
		},
		"tqdm": {
			Name:         "tqdm",
			Version:      "4.64.0",
			Dependencies: []string{},
		},
	}

	assert.Equal(t, gt, packageMetadata)
	assert.Nil(t, job.err)
}
