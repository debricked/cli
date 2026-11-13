package bundler

import (
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/sanitize/pm/util"
)

const (
	gemSection = "GEM"
	gitSection = "GIT"
	specsLine  = "  specs:"
	specIndent = "    "
)

var (
	// gemSectionRegex matches section headers such as GEM, PLATFORMS and RUBY VERSION.
	gemSectionRegex = regexp.MustCompile(`^[A-Z][A-Z ]*$`)
	// gemSpecRegex matches a package line, `    name (version)`. The version can name a platform.
	gemSpecRegex = regexp.MustCompile(`^ {4}[^\s()]+ \([^\s()]+\)$`)
)

// sanitizeGemfileLock keeps the package lines of GEM blocks and the first package line of each
// GIT block, the only lines the backend reads. remote:, revision: and branch: lines are removed,
// because they can hold credentials and branch names. PATH blocks hold local paths and the backend
// ignores them, so they are removed with every other section. Blocks are separated by a blank
// line, which is how the backend splits them.
func sanitizeGemfileLock(content []byte, _ *util.FileContext) ([]byte, error) {
	sanitizer := &gemfileLockSanitizer{}
	for _, line := range util.SplitLines(content) {
		if err := sanitizer.processLine(line); err != nil {
			return nil, err
		}
	}
	sanitizer.endSection()
	if sanitizer.section == "" {
		return nil, util.ErrInvalidFormat
	}

	return []byte(strings.Join(sanitizer.blocks, "\n")), nil
}

type gemfileLockSanitizer struct {
	blocks  []string
	block   *gemBlock
	section string
	inSpecs bool
}

func (s *gemfileLockSanitizer) processLine(line string) error {
	switch {
	case strings.TrimSpace(line) == "":
		return nil
	case line[0] != ' ':
		return s.startSection(line)
	case s.section == "":
		return util.ErrInvalidFormat
	case s.block == nil:
		return nil
	case line == specsLine:
		s.inSpecs = true

		return nil
	case s.inSpecs && strings.HasPrefix(line, specIndent) && line[len(specIndent)] != ' ':
		return s.addSpec(line)
	default:
		return nil
	}
}

func (s *gemfileLockSanitizer) startSection(line string) error {
	if !gemSectionRegex.MatchString(line) {
		return util.ErrInvalidFormat
	}
	s.endSection()
	s.section, s.inSpecs = line, false
	if line == gemSection || line == gitSection {
		s.block = &gemBlock{header: line}
	}

	return nil
}

// addSpec keeps every package of a GEM block, and only the first of a GIT block.
func (s *gemfileLockSanitizer) addSpec(line string) error {
	if !gemSpecRegex.MatchString(line) {
		return util.ErrInvalidFormat
	}
	if s.section == gemSection || len(s.block.specs) == 0 {
		s.block.specs = append(s.block.specs, line)
	}

	return nil
}

func (s *gemfileLockSanitizer) endSection() {
	if s.block != nil && len(s.block.specs) > 0 {
		s.blocks = append(s.blocks, s.block.String())
	}
	s.block = nil
}

type gemBlock struct {
	header string
	specs  []string
}

func (b *gemBlock) String() string {
	return b.header + "\n" + specsLine + "\n" + strings.Join(b.specs, "\n") + "\n"
}
