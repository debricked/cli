package file

import (
	"path"
	"regexp"

	"github.com/debricked/cli/internal/resolution/pm"
	"github.com/debricked/cli/internal/resolution/pm/npm"
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
		for _, p := range bf.pms {
			if bf.skipPackageManager(p) {
				continue
			}

			for _, manifest := range p.Manifests() {
				compiledRegex, _ := regexp.Compile(manifest)
				if compiledRegex.MatchString(path.Base(file)) {
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
