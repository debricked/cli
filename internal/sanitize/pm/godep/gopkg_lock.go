package godep

import (
	"strings"

	"github.com/debricked/cli/internal/sanitize/pm/util"
	"github.com/pelletier/go-toml/v2"
)

type gopkgLock struct {
	Projects []struct {
		Name    string `toml:"name"`
		Version string `toml:"version"`
	} `toml:"projects"`
}

// sanitizeGopkgLock keeps the name and version of each project. A name can contain a private
// host, but it is what identifies the dependency, so it is kept. Branches, revisions, sources,
// import paths and [solve-meta], which lists the project's own imports, are removed.
func sanitizeGopkgLock(content []byte, _ *util.FileContext) ([]byte, error) {
	var lock gopkgLock
	if err := toml.Unmarshal(content, &lock); err != nil || len(lock.Projects) == 0 {
		return nil, util.ErrInvalidFormat
	}

	blocks := make([]string, 0, len(lock.Projects))
	for _, project := range lock.Projects {
		if project.Name == "" {
			return nil, util.ErrInvalidFormat
		}
		block := "[[projects]]\n  name = " + util.QuoteTOML(project.Name) + "\n"
		if project.Version != "" {
			block += "  version = " + util.QuoteTOML(project.Version) + "\n"
		}
		blocks = append(blocks, block)
	}

	return []byte(strings.Join(blocks, "\n")), nil
}
