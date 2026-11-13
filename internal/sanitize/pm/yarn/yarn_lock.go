package yarn

import (
	"bytes"
	"net/url"
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/sanitize/pm/util"
)

const (
	yarnMetadataHeader  = "__metadata:"
	yarnProjectSpec     = "workspace:."
	yarnWorkspacePrefix = "workspace:"
	yarnBindingsMarker  = "::"
	yarnLocatorBinding  = "locator"
	yarnSpecSeparator   = ", "
	yarnKeyIndent       = 2
	yarnDependencyLevel = 4
)

// yarnEntryRegex matches both `version "1.0.0"` (v1) and `version: 1.0.0` (Berry).
var yarnEntryRegex = regexp.MustCompile(`^("[^"]+"|[^\s:"]+)(:?\s+)(.*)$`)

// yarnLockSanitizer filters the file line by line rather than regenerating it, because the backend
// parser also works line by line and depends on the blank lines and the indentation.
type yarnLockSanitizer struct {
	context      *util.FileContext
	projectNames map[string]bool
	// projectDescriptors holds every `name@range` in the header of a project's own block. A
	// registry package can have the project's name, so a name alone doesn't identify the project.
	projectDescriptors map[string]bool
	output             bytes.Buffer
	inBlock            bool
	inMetadata         bool
	inDependency       bool
}

func sanitizeYarnLock(content []byte, context *util.FileContext) ([]byte, error) {
	sanitizer := &yarnLockSanitizer{
		context:            context,
		projectNames:       map[string]bool{},
		projectDescriptors: map[string]bool{},
	}
	lines := strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n")

	if err := sanitizer.findProjectNames(lines); err != nil {
		return nil, err
	}
	for _, line := range lines {
		if err := sanitizer.processLine(line); err != nil {
			return nil, err
		}
	}

	return sanitizer.output.Bytes(), nil
}

// findProjectNames runs first, so that a dependency on the project is renamed even when it appears
// before the project's own block.
func (s *yarnLockSanitizer) findProjectNames(lines []string) error {
	for _, line := range lines {
		if !isYarnHeader(line) || line == yarnMetadataHeader {
			continue
		}
		if err := s.findProjectDescriptors(line); err != nil {
			return err
		}
	}

	return nil
}

// findProjectDescriptors records the specs of a header that has a `<name>@workspace:.` spec. All
// specs of that name in the header resolve to the project, such as `<name>@workspace:*`.
func (s *yarnLockSanitizer) findProjectDescriptors(line string) error {
	specs := strings.Split(strings.TrimSuffix(line, ":"), yarnSpecSeparator)
	projectName := ""
	for _, spec := range specs {
		name, version, valid := splitYarnSpec(strings.Trim(spec, `"`))
		if !valid {
			return util.ErrInvalidFormat
		}
		if version == yarnProjectSpec {
			projectName = name
		}
	}
	if projectName == "" {
		return nil
	}
	s.projectNames[projectName] = true
	for _, spec := range specs {
		unquoted := strings.Trim(spec, `"`)
		if name, _, _ := splitYarnSpec(unquoted); name == projectName {
			s.projectDescriptors[unquoted] = true
		}
	}

	return nil
}

func isYarnHeader(line string) bool {
	return line != "" && line[0] != ' ' && line[0] != '#'
}

func (s *yarnLockSanitizer) processLine(line string) error {
	trimmed := strings.TrimLeft(line, " ")
	if trimmed == "" || trimmed[0] == '#' {
		return nil
	}
	indent := len(line) - len(trimmed)

	switch {
	case indent == 0:
		return s.processHeader(line)
	case !s.inBlock:
		return util.ErrInvalidFormat
	case s.inMetadata:
		return nil
	case indent == yarnKeyIndent:
		return s.processKey(trimmed)
	case indent >= yarnDependencyLevel:
		if indent == yarnDependencyLevel && s.inDependency {
			return s.processDependency(trimmed)
		}

		return nil
	}

	return util.ErrInvalidFormat
}

func (s *yarnLockSanitizer) processHeader(line string) error {
	if !strings.HasSuffix(line, ":") {
		return util.ErrInvalidFormat
	}
	s.inBlock = true
	s.inDependency = false
	s.inMetadata = line == yarnMetadataHeader
	if s.inMetadata {
		return nil
	}

	specs := strings.Split(strings.TrimSuffix(line, ":"), yarnSpecSeparator)
	for index, spec := range specs {
		masked, err := s.maskHeaderSpec(spec)
		if err != nil {
			return err
		}
		specs[index] = masked
	}

	if s.output.Len() > 0 {
		s.output.WriteByte('\n')
	}
	s.writeLine(strings.Join(specs, yarnSpecSeparator) + ":")

	return nil
}

// maskHeaderSpec keeps any surrounding quote where it was. In Berry one pair of quotes surrounds
// the whole header, so a spec can have a quote on only one side.
func (s *yarnLockSanitizer) maskHeaderSpec(spec string) (string, error) {
	inner := strings.TrimLeft(spec, `"`)
	leading := spec[:len(spec)-len(inner)]
	unquoted := strings.TrimRight(inner, `"`)
	trailing := inner[len(unquoted):]

	name, version, valid := splitYarnSpec(unquoted)
	if !valid {
		return "", util.ErrInvalidFormat
	}

	return leading + s.projectName(name, version) + "@" + s.maskHeaderVersion(version) + trailing, nil
}

// maskHeaderVersion keeps Berry's `::locator=` suffix, in redacted form. Berry adds it to local
// dependencies such as portal: and link:, and the dependency lists refer to them without it. The
// backend needs that difference to give the same result as the original, and the suffix keeps
// the headers of the same path from different workspaces apart once the paths are masked.
func (s *yarnLockSanitizer) maskHeaderVersion(version string) string {
	rangePart, bindings, found := strings.Cut(version, yarnBindingsMarker)
	if !found {
		return util.MaskValue(version)
	}
	parsed, err := url.ParseQuery(bindings)
	locator := parsed.Get(yarnLocatorBinding)
	if err != nil || locator == "" {
		return util.MaskValue(version)
	}
	name, locatorRange, valid := splitYarnSpec(locator)
	if !valid {
		return util.MaskValue(version)
	}
	masked := s.projectName(name, locatorRange) + "@" + util.MaskValue(locatorRange)

	return util.MaskValue(rangePart) + yarnBindingsMarker + yarnLocatorBinding + "=" + url.QueryEscape(masked)
}

// splitYarnSpec splits at the first '@' after the start, so scoped names and values that contain
// '@' themselves, such as git URLs, stay whole.
func splitYarnSpec(spec string) (string, string, bool) {
	if len(spec) < 2 {
		return "", "", false
	}
	index := strings.Index(spec[1:], "@")
	if index < 0 {
		return "", "", false
	}

	return spec[:index+1], spec[index+2:], true
}

func (s *yarnLockSanitizer) processKey(entry string) error {
	s.inDependency = false

	if entry == "dependencies:" || entry == "optionalDependencies:" {
		s.inDependency = true
		s.writeLine("  " + entry)

		return nil
	}

	match := yarnEntryRegex.FindStringSubmatch(entry)
	if match == nil {
		return nil
	}
	key, separator, value := match[1], match[2], match[3]
	unquoted := strings.Trim(value, `"`)

	switch {
	case key == "version":
		s.writeLine("  " + key + separator + maskQuoted(value))
	case key == "resolved" && util.MatchesUnsupportedProtocol(unquoted):
		s.writeLine("  " + key + separator + maskQuoted(value))
	}

	return nil
}

func (s *yarnLockSanitizer) processDependency(entry string) error {
	match := yarnEntryRegex.FindStringSubmatch(entry)
	if match == nil {
		return util.ErrInvalidFormat
	}
	key, separator, value := match[1], match[2], match[3]

	name := strings.Trim(key, `"`)
	renamed := strings.Replace(key, name, s.projectName(name, strings.Trim(value, `"`)), 1)
	s.writeLine("    " + renamed + separator + maskQuoted(value))

	return nil
}

// projectName renames only a descriptor that refers to the project. A package from a registry keeps
// its name even when the project has the same name, such as a tool that formats itself with a
// published copy of itself.
func (s *yarnLockSanitizer) projectName(name, version string) string {
	if !s.projectNames[name] {
		return name
	}
	if s.projectDescriptors[name+"@"+version] || strings.HasPrefix(version, yarnWorkspacePrefix) {
		return s.context.ProjectPlaceholder(name)
	}

	return name
}

func (s *yarnLockSanitizer) writeLine(line string) {
	s.output.WriteString(line)
	s.output.WriteByte('\n')
}

func maskQuoted(value string) string {
	unquoted := strings.Trim(value, `"`)

	return strings.Replace(value, unquoted, util.MaskValue(unquoted), 1)
}
