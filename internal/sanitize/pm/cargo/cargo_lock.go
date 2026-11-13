package cargo

import (
	"strings"

	"github.com/debricked/cli/internal/sanitize/pm/util"
	"github.com/pelletier/go-toml/v2"
)

const (
	maskedGitSource = "git+redacted"
	// localCrateVersion replaces the version of the customer's own crates, which are renamed.
	localCrateVersion = "0.0.0"
)

type cargoLock struct {
	Packages []struct {
		Name    string `toml:"name"`
		Version string `toml:"version"`
		Source  string `toml:"source"`
	} `toml:"package"`
}

// sanitizeCargoLock keeps the name and version of each package. The backend only reads whether a
// source starts with "git", so a git source is replaced and other sources are removed. A package
// with no source is the customer's own crate or a path dependency, so it is renamed.
//
// The backend reads the lines after each name without checking where the file ends, so every
// block, including the last, ends with a blank line.
func sanitizeCargoLock(content []byte, context *util.FileContext) ([]byte, error) {
	var lock cargoLock
	if err := toml.Unmarshal(content, &lock); err != nil || len(lock.Packages) == 0 {
		return nil, util.ErrInvalidFormat
	}

	var output strings.Builder
	for _, crate := range lock.Packages {
		if crate.Name == "" {
			return nil, util.ErrInvalidFormat
		}
		name, version := crate.Name, crate.Version
		if crate.Source == "" {
			name, version = context.ProjectPlaceholder(crate.Name), localCrateVersion
		}
		output.WriteString("[[package]]\n")
		output.WriteString("name = " + util.QuoteTOML(name) + "\n")
		output.WriteString("version = " + util.QuoteTOML(version) + "\n")
		if strings.HasPrefix(crate.Source, "git") {
			output.WriteString("source = " + util.QuoteTOML(maskedGitSource) + "\n")
		}
		output.WriteString("\n")
	}

	return []byte(output.String()), nil
}
