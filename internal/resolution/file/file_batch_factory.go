package file

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/pm"
	"github.com/debricked/cli/internal/resolution/pm/npm"
	"github.com/debricked/cli/internal/resolution/pm/poetry"
	"github.com/debricked/cli/internal/resolution/pm/uv"
	"github.com/debricked/cli/internal/resolution/pm/yarn"
)

type IBatchFactory interface {
	Make(files []string) []IBatch
	SetNpmPreferred(npmPreferred bool)
}

type BatchFactory struct {
	pms          []pm.IPm
	npmPreferred bool
}

func NewBatchFactory() *BatchFactory {
	return &BatchFactory{
		pms: pm.Pms(),
	}
}

func (bf *BatchFactory) SetNpmPreferred(npmPreferred bool) {
	bf.npmPreferred = npmPreferred
}

func (bf *BatchFactory) Make(files []string) []IBatch {
	batchMap := make(map[string]IBatch)
	for _, file := range files {
		base := filepath.Base(file)
		for _, p := range bf.pms {
			if bf.skipPackageManager(p) {
				continue
			}

			for _, manifest := range p.Manifests() {
				// Special handling for Python pyproject.toml, which may belong to either Poetry or uv.
				if manifest == "pyproject.toml" && strings.EqualFold(base, "pyproject.toml") {
					pmName := detectPyprojectPm(file)
					if pmName != p.Name() {
						continue
					}
				}

				compiledRegex, _ := regexp.Compile(manifest)
				if compiledRegex.MatchString(base) {
					batch, ok := batchMap[p.Name()]
					if !ok {
						batch = NewBatch(p)
						batchMap[p.Name()] = batch
					}
					batch.Add(file)
				}
			}
		}
	}

	batches := make([]IBatch, 0, len(batchMap))

	for _, batch := range batchMap {
		batches = append(batches, batch)
	}

	return batches
}

func (bf *BatchFactory) skipPackageManager(p pm.IPm) bool {
	name := p.Name()

	switch true {
	case name == npm.Name && !bf.npmPreferred:
		return true
	case name == yarn.Name && bf.npmPreferred:
		return true
	}

	return false
}

func detectPyprojectPm(pyprojectPath string) string {
	dir := filepath.Dir(pyprojectPath)

	if fileExists(filepath.Join(dir, "uv.lock")) {
		return uv.Name
	}

	if fileExists(filepath.Join(dir, "poetry.lock")) {
		return poetry.Name
	}

	content, err := os.ReadFile(pyprojectPath)
	if err == nil {
		data := string(content)
		hasPoetry := strings.Contains(data, "[tool.poetry]") ||
			strings.Contains(data, "tool.poetry")
		hasProject := strings.Contains(data, "[project]")

		if hasPoetry && hasProject {
			// Ambiguous: both Poetry and UV indicators present
			return ""
		}
		if hasPoetry {
			return poetry.Name
		}
		if hasProject {
			return uv.Name
		}
	}

	// If no indicators found, cannot determine PM
	return ""
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !info.IsDir()
}
