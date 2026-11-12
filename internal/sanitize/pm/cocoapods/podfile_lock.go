package cocoapods

import (
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/sanitize/pm/util"
)

const (
	podsSection            = "PODS:"
	externalSourcesSection = "EXTERNAL SOURCES:"
)

var (
	// podSectionRegex matches section lines such as PODS: and PODFILE CHECKSUM: cccc.
	podSectionRegex = regexp.MustCompile(`^[A-Z][A-Z ]*:`)
	// podRegex matches `  - Name (version)`, quoted or not, with the colon that starts a pod's
	// own dependency list removed first.
	podRegex = regexp.MustCompile(`^  - (?:"[^"\s()]+ \([^"()]+\)"|[^"\s()]+ \([^"()]+\))$`)
	// externalSourceRegex matches `  Name:`, the line that names a pod from a git repository or
	// a local path.
	externalSourceRegex = regexp.MustCompile(`^  (?:"[^"]+"|[^"\s]+):$`)
)

// sanitizePodfileLock keeps the PODS lines and the names under EXTERNAL SOURCES, which decide the
// pods that have no version. Each pod's dependencies, its git URL, path, branch and commit, and
// every other section are removed. Each section ends with a blank line, because the backend finds
// a section by the blank line after it.
func sanitizePodfileLock(content []byte, _ *util.FileContext) ([]byte, error) {
	sanitizer := &podfileLockSanitizer{sections: map[string][]string{}}
	for _, line := range util.SplitLines(content) {
		if err := sanitizer.processLine(line); err != nil {
			return nil, err
		}
	}
	if _, found := sanitizer.sections[podsSection]; !found {
		return nil, util.ErrInvalidFormat
	}

	var output strings.Builder
	for _, name := range []string{podsSection, externalSourcesSection} {
		lines, found := sanitizer.sections[name]
		if !found {
			continue
		}
		output.WriteString(name + "\n")
		for _, line := range lines {
			output.WriteString(line + "\n")
		}
		output.WriteString("\n")
	}

	return []byte(output.String()), nil
}

type podfileLockSanitizer struct {
	sections map[string][]string
	section  string
}

func (s *podfileLockSanitizer) processLine(line string) error {
	switch {
	case strings.TrimSpace(line) == "":
		return nil
	case line[0] != ' ':
		return s.startSection(line)
	case s.section == "":
		return util.ErrInvalidFormat
	case s.section == podsSection && strings.HasPrefix(line, "  - "):
		return s.keep(strings.TrimSuffix(line, ":"), podRegex)
	case s.section == externalSourcesSection && strings.HasPrefix(line, "  ") && line[2] != ' ':
		return s.keep(line, externalSourceRegex)
	default:
		return nil
	}
}

func (s *podfileLockSanitizer) startSection(line string) error {
	if !podSectionRegex.MatchString(line) {
		return util.ErrInvalidFormat
	}
	s.section = line
	if line == podsSection || line == externalSourcesSection {
		s.sections[line] = []string{}
	}

	return nil
}

func (s *podfileLockSanitizer) keep(line string, allowed *regexp.Regexp) error {
	if !allowed.MatchString(line) {
		return util.ErrInvalidFormat
	}
	s.sections[s.section] = append(s.sections[s.section], line)

	return nil
}
