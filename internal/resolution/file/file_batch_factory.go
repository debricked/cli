package file

import (
	"path"
	"regexp"

	"github.com/debricked/cli/internal/resolution/pm"
)

type IBatchFactory interface {
	Make(files []string) []IBatch
}

type BatchFactory struct {
	pms []pm.IPm
}

func NewBatchFactory() BatchFactory {
	return BatchFactory{
		pms: pm.Pms(),
	}
}

func (bf BatchFactory) Make(files []string) []IBatch {
	batchMap := make(map[string]IBatch)
	for _, file := range files {
		for _, p := range bf.pms {
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
