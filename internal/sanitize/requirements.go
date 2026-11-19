package sanitize

import (
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/sanitize/pm/util"
)

var (
	requirementVersionComment = regexp.MustCompile(`^#version=[A-Za-z0-9][A-Za-z0-9._-]*$`)
	requirementLine           = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._-]*(?:\[[A-Za-z0-9._,-]+\])?(?:\s*(?:===|==|!=|~=|<=|>=|<|>)\s*[A-Za-z0-9][A-Za-z0-9.*+!_-]*(?:\s*,\s*(?:===|==|!=|~=|<=|>=|<|>)\s*[A-Za-z0-9][A-Za-z0-9.*+!_-]*)*)?(?:\s*;\s*[A-Za-z0-9_.' "<>=!()]+)?$`)
)

func sanitizeRequirements(content []byte, _ *util.FileContext) ([]byte, error) {
	for _, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			if requirementVersionComment.MatchString(line) {
				continue
			}

			return nil, util.ErrInvalidFormat
		}
		if !requirementLine.MatchString(line) {
			return nil, util.ErrInvalidFormat
		}
	}

	return content, nil
}

func (r *run) copyRequirementsFiles(exclusions, inclusions []string) error {
	return filepath.WalkDir(r.root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || filepath.Base(path) != "requirements.txt" || file.Excluded(exclusions, inclusions, path) {
			return nil
		}
		r.processFile(path, sanitizeRequirements, nil)

		return nil
	})
}
