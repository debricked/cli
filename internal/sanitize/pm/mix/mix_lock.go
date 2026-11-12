package mix

import (
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/sanitize/pm/util"
)

const (
	mixOpen  = "%{"
	mixClose = "}"
)

var (
	mixEntryRegex = regexp.MustCompile(`^  "[^"]+": \{`)
	mixHexRegex   = regexp.MustCompile(`^  ("[^"]+"): \{:hex, (:\w+), ("[^"\s]+")[,}]`)
	mixGitRegex   = regexp.MustCompile(`^  ("[^"]+"): \{:git[,}]`)
)

// sanitizeMixLock writes one line per package, as the backend reads them with a regex: name,
// product and version for Hex packages, and the name alone for git packages, which have no
// version. Git URLs, commits, branches, hashes and repository names such as hexpm:acme_org are
// removed. Entries of other kinds are removed too, because the backend doesn't read them.
func sanitizeMixLock(content []byte, _ *util.FileContext) ([]byte, error) {
	var lines []string
	for _, line := range util.SplitLines(content) {
		if strings.TrimSpace(line) != "" {
			lines = append(lines, line)
		}
	}
	if len(lines) < 2 || lines[0] != mixOpen || lines[len(lines)-1] != mixClose {
		return nil, util.ErrInvalidFormat
	}

	var output strings.Builder
	output.WriteString(mixOpen + "\n")
	for _, line := range lines[1 : len(lines)-1] {
		entry, err := sanitizeMixEntry(line)
		if err != nil {
			return nil, err
		}
		output.WriteString(entry)
	}
	output.WriteString(mixClose + "\n")

	return []byte(output.String()), nil
}

// sanitizeMixEntry returns the sanitized line for one entry, or "" for an entry of another kind.
func sanitizeMixEntry(line string) (string, error) {
	if !mixEntryRegex.MatchString(line) {
		return "", util.ErrInvalidFormat
	}
	if hex := mixHexRegex.FindStringSubmatch(line); hex != nil {
		return "  " + hex[1] + ": {:hex, " + hex[2] + ", " + hex[3] + "},\n", nil
	}
	if strings.Contains(line, "{:hex") {
		return "", util.ErrInvalidFormat
	}
	if git := mixGitRegex.FindStringSubmatch(line); git != nil {
		return "  " + git[1] + ": {:git},\n", nil
	}

	return "", nil
}
