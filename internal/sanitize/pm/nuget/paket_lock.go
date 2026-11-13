package nuget

import (
	"bytes"
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/sanitize/pm/util"
)

// paketPackageRegex matches the lines the backend parser reads: exactly four spaces, then
// `name (version)`. Anything after the closing bracket, such as a restriction, is left out.
var paketPackageRegex = regexp.MustCompile(`^ {4}(\S.*?) \(([^)]*)\)`)

// paketHeaders holds the top-level words of a paket.lock, used to tell a real one from other text.
var paketHeaders = map[string]bool{
	"REDIRECTS:": true, "RESTRICTION:": true, "STORAGE:": true, "STRATEGY:": true, "LOWEST_MATCHING:": true,
	"FRAMEWORK:": true, "CONTENT:": true, "COPY-CONTENT-TO-OUTPUT-DIR:": true, "COPY-LOCAL:": true,
	"IMPORT-TARGETS:": true, "CLITOOL:": true, "VERSION:": true, "AUTO-DETECT:": true, "GENERATE-LOAD-SCRIPTS:": true,
	"GROUP": true, "NUGET": true, "GITHUB": true, "GIT": true, "HTTP": true,
}

// sanitizePaketLock outputs nothing but `    name (version)` lines. GIT and HTTP entries don't
// match the parser, and GITHUB entries are renamed because the backend reports their source paths
// as packages.
func sanitizePaketLock(content []byte, context *util.FileContext) ([]byte, error) {
	var output bytes.Buffer
	section := ""

	for _, line := range strings.Split(strings.ReplaceAll(string(content), "\r\n", "\n"), "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if line[0] != ' ' {
			fields := strings.Fields(line)
			if !paketHeaders[fields[0]] {
				return nil, util.ErrInvalidFormat
			}
			section = fields[0]

			continue
		}

		match := paketPackageRegex.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		name, version := match[1], match[2]

		switch section {
		case "NUGET":
			output.WriteString("    " + name + " (" + version + ")\n")
		case "GITHUB":
			output.WriteString("    " + context.FilePlaceholder(name) + " (" + version + ")\n")
		}
	}

	return output.Bytes(), nil
}
