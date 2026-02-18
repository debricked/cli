package file

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/pm"
	"github.com/debricked/cli/internal/resolution/pm/npm"
	"github.com/debricked/cli/internal/resolution/pm/pnpm"
	"github.com/debricked/cli/internal/resolution/pm/poetry"
	"github.com/debricked/cli/internal/resolution/pm/uv"
	"github.com/debricked/cli/internal/resolution/pm/yarn"
	"github.com/fatih/color"
)

type IBatchFactory interface {
	Make(files []string) []IBatch
	SetNpmPreferred(npmPreferred bool)
}

type BatchFactory struct {
	pms                 []pm.IPm
	npmPreferred        bool
	warnedYarnDefaultPM bool
}

func NewBatchFactory() *BatchFactory {
	return &BatchFactory{
		pms: pm.Pms(),
	}
}

func (bf *BatchFactory) SetNpmPreferred(npmPreferred bool) {
	bf.npmPreferred = npmPreferred
}

//nolint:cyclop
func (bf *BatchFactory) Make(files []string) []IBatch {
	batchMap := make(map[string]IBatch)
	for _, file := range files {
		bf.processFile(file, batchMap)
	}

	batches := make([]IBatch, 0, len(batchMap))
	for _, batch := range batchMap {
		batches = append(batches, batch)
	}

	return batches
}

func (bf *BatchFactory) processFile(file string, batchMap map[string]IBatch) {
	base := filepath.Base(file)
	for _, p := range bf.pms {
		for _, manifest := range p.Manifests() {
			if bf.shouldProcessManifest(manifest, base, file, p) {
				compiledRegex, _ := regexp.Compile(manifest)
				if compiledRegex.MatchString(base) {
					bf.addToBatch(p, file, batchMap)
				}
			}
		}
	}
}

func (bf *BatchFactory) shouldProcessManifest(manifest, base, file string, p pm.IPm) bool {
	if isNodePackageJSON(manifest, base) {
		return bf.shouldProcessNodeManifest(file, p)
	}

	if isPyprojectToml(manifest, base) {
		return shouldProcessPyprojectManifest(file, p)
	}

	return true
}

func isNodePackageJSON(manifest, base string) bool {
	return manifest == `package\.json$` && strings.EqualFold(base, "package.json")
}

func (bf *BatchFactory) shouldProcessNodeManifest(file string, p pm.IPm) bool {
	pmName := detectNodePm(file)
	if pmName != "" {
		// If we can detect the PM from lockfiles or package.json, use that
		return pmName == p.Name()
	}

	// No explicit packageManager found: fall back to npmPreferred flag between npm and yarn
	switch {
	case p.Name() == npm.Name && bf.npmPreferred:
		return true
	case p.Name() == yarn.Name && !bf.npmPreferred:
		bf.warnYarnDefault()

		return true
	default:
		return false
	}
}

func isPyprojectToml(manifest, base string) bool {
	return manifest == "pyproject.toml" && strings.EqualFold(base, "pyproject.toml")
}

func shouldProcessPyprojectManifest(file string, p pm.IPm) bool {
	pmName := detectPyprojectPm(file)
	return pmName == p.Name()
}

func (bf *BatchFactory) addToBatch(p pm.IPm, file string, batchMap map[string]IBatch) {
	batch, ok := batchMap[p.Name()]
	if !ok {
		batch = NewBatch(p)
		batchMap[p.Name()] = batch
	}
	batch.Add(file)
}

func (bf *BatchFactory) warnYarnDefault() {
	if bf.warnedYarnDefaultPM {

		return
	}

	fmt.Printf("%s  Unable to detect package manager through package.json file, defaulting to yarn.\n", color.YellowString("⚠️"))
	bf.warnedYarnDefaultPM = true
}

func detectNodePm(packageJSONPath string) string {
	return detectNodePmFromPackageJSON(packageJSONPath)
}

func detectNodePmFromPackageJSON(packageJSONPath string) string {
	// Prefer explicit packageManager field if present
	content, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return ""
	}

	var pkg struct {
		PackageManager string `json:"packageManager"`
	}
	if jsonErr := json.Unmarshal(content, &pkg); jsonErr != nil || pkg.PackageManager == "" {
		return ""
	}

	name := pkg.PackageManager
	if at := strings.Index(name, "@"); at > 0 {
		name = name[:at]
	}

	switch name {
	case pnpm.Name:
		return pnpm.Name
	case yarn.Name:
		return yarn.Name
	case npm.Name:
		return npm.Name
	default:
		return ""
	}
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
