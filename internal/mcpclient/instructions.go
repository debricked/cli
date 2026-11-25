package mcpclient

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/debricked/cli/internal/io"
)

const instructionStart = "<!-- debricked:dependency-policy:start -->"
const instructionEnd = "<!-- debricked:dependency-policy:end -->"

const dependencyPolicy = "## Dependency policy\n\n" +
	"Before adding, installing, upgrading, importing, or recommending a dependency, call\n" +
	"`check_dependency_policy_compliance` with 1–5 distinct, versioned PURLs (`purls`); use one for a\n" +
	"single package.\n\n" +
	"- **Compliant:** Report the checked package/version as approved by Fortify SCA (Debricked).\n" +
	"- **Blocked:** Do not add it. Report every rule's condition, package-specific trigger evidence, and\n" +
	"  configuration URL when present. Disclose truncated results; never infer which OR branch matched.\n" +
	"- **Unavailable:** Report the reason; do not claim approval.\n" +
	"- **Alternatives:** For a blocked package, inspect its use and compatibility, then present up to\n" +
	"  five exact, versioned PURLs as unchecked candidates. Ask approval for that exact list. After\n" +
	"  approval, check all approved candidates in one call. Recommend only candidates that pass; report\n" +
	"  each result and say if none pass. Approval permits checking only, not adding or installing.\n" +
	"  Get fresh approval before checking a different candidate list."

// InstructionsPath always resolves project instructions, independently of the MCP config scope.
func InstructionsPath(client Client, env IEnvironment) (string, error) {
	wd, err := env.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to resolve current directory: %w", err)
	}

	switch client {
	case ClientVSCode:
		return filepath.Join(wd, ".github", "copilot-instructions.md"), nil
	case ClientClaude:
		return filepath.Join(wd, "CLAUDE.md"), nil
	case ClientCursor:
		return filepath.Join(wd, ".cursor", "rules", "debricked.mdc"), nil
	default:
		return "", fmt.Errorf("unsupported client %q", client)
	}
}

// SetupInstructions refreshes only the managed policy block in place.
func SetupInstructions(fs io.IFileSystem, path string, client Client) (changed, existed bool, err error) {
	raw, err := fs.ReadFile(path)
	if err != nil && !fs.IsNotExist(err) {
		return false, false, fmt.Errorf("failed to read %s: %w", path, err)
	}
	existed = err == nil

	content, err := upsertInstructions(string(raw), client)
	if err != nil {
		return false, existed, fmt.Errorf("failed to update %s: %w", path, err)
	}
	if content == string(raw) {
		return false, existed, nil
	}

	if err := fs.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return false, existed, fmt.Errorf("failed to create directory for %s: %w", path, err)
	}
	if err := fs.FsWriteFile(path, []byte(content), 0644); err != nil {
		return false, existed, fmt.Errorf("failed to write %s: %w", path, err)
	}

	return true, existed, nil
}

func upsertInstructions(content string, client Client) (string, error) {
	newline := "\n"
	if strings.Contains(content, "\r\n") {
		newline = "\r\n"
	}
	block := strings.ReplaceAll(instructionStart+"\n"+dependencyPolicy+"\n"+instructionEnd, "\n", newline)
	start, end := strings.Index(content, instructionStart), strings.Index(content, instructionEnd)
	if start >= 0 || end >= 0 {
		if start < 0 || end < start || strings.Count(content, instructionStart) != 1 || strings.Count(content, instructionEnd) != 1 {
			return "", fmt.Errorf("invalid dependency policy markers; fix or remove them manually")
		}
		content = content[:start] + block + content[end+len(instructionEnd):]
	} else {
		if content != "" {
			if !strings.HasSuffix(content, newline) {
				content += newline
			}
			content += newline
		}
		content += block + newline
	}

	if client == ClientCursor {
		return alwaysApplyCursorRule(content, newline)
	}

	return content, nil
}

func alwaysApplyCursorRule(content, newline string) (string, error) {
	delimiter := "---" + newline
	if !strings.HasPrefix(content, delimiter) {
		return delimiter + "alwaysApply: true" + newline + delimiter + newline + content, nil
	}

	lines := strings.Split(content, newline)
	for i := 1; i < len(lines); i++ {
		if lines[i] == "---" {
			header := make([]string, 0, i+1)
			header = append(header, lines[0])
			for _, line := range lines[1:i] {
				if !strings.HasPrefix(line, "alwaysApply:") {
					header = append(header, line)
				}
			}
			header = append(header, "alwaysApply: true")

			return strings.Join(append(header, lines[i:]...), newline), nil
		}
	}

	return "", fmt.Errorf("unterminated Cursor rule frontmatter; fix it manually")
}
