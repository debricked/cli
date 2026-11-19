package util

import (
	"net/url"
	"regexp"
	"strings"
)

const (
	MaskedGitWord         = "git+ssh://git@redacted/redacted.git#redacted"
	MaskedGitOther        = "git+ssh://git@redacted/redacted-repo.git#redacted"
	MaskedGitWordNoUser   = "git+ssh://redacted/redacted.git#redacted"
	MaskedGitOtherNoUser  = "git+ssh://redacted/redacted-repo.git#redacted"
	MaskedLocalDot        = "file:./redacted"
	MaskedLocal           = "file:redacted"
	MaskedURL             = "https://redacted.invalid/redacted.tgz"
	MaskedWorkspace       = "workspace:redacted"
	MaskedPatchFile       = "patch:redacted.patch"
	MaskedPatch           = "patch:redacted"
	MaskedUnknown         = "redacted"
	yarnBuiltinPatchToken = "optional!builtin"
	patchProtocol         = "patch:"
	maskedPatchFileName   = "redacted.patch"
	maskedPatchSource     = "redacted"
)

var (
	// plainVersionRegex matches version ranges and dist-tags. It excludes ':', '/' and '@',
	// so that nothing shaped like a URL, a path or a host can pass as a version.
	plainVersionRegex = regexp.MustCompile(`^[\w\s.*+^~<>=|,!-]*$`)
	npmProtocolRegex  = regexp.MustCompile(`^npm:(?:@?[\w.\-]+(?:/[\w.\-]+)?@)?[\w\s.*+^~<>=|,!-]*$`)
	gitShorthandRegex = regexp.MustCompile(`^[\w.\-]+/[\w.\-]+(?:#.*)?$`)
	wordOnlyRegex     = regexp.MustCompile(`^\w+$`)
	patchLocatorRegex = regexp.MustCompile(`&?locator=[^&]*`)
	// packageNameRegex matches an npm package name, scoped or not, and nothing shaped like a path.
	packageNameRegex = regexp.MustCompile(`^(?:@[\w.\-]+/)?[\w.\-]+$`)
)

var publicRegistryHosts = map[string]bool{
	"registry.npmjs.org":   true,
	"registry.yarnpkg.com": true,
}

var localProtocols = []string{"file:", "link:", "portal:", "exec:"}

var gitPrefixes = []string{"git+", "git://", "ssh://", "git@", "github:", "gitlab:", "bitbucket:", "gist:"}

// IsPlainVersion reports whether value is a version, a version range or a dist-tag, with nothing
// that could name a host or a path.
func IsPlainVersion(value string) bool {
	return plainVersionRegex.MatchString(value)
}

// MaskValue replaces any private source in a dependency value with a placeholder that the
// backend classifies the same way as the original. The function is pure, so a value gets the same
// placeholder in a lock file and in its manifest.
//
// Anything it doesn't recognise is masked rather than kept, so a new value format can't leak.
func MaskValue(value string) string {
	switch {
	case plainVersionRegex.MatchString(value), npmProtocolRegex.MatchString(value):
		return value
	case isPublicRegistryURL(value):
		return value
	case strings.HasPrefix(value, "workspace:"):
		return maskWorkspace(value)
	case strings.HasPrefix(value, "patch:"):
		return maskPatch(value)
	}

	if rest, isLocal := CutLocalProtocol(value); isLocal {
		return maskLocal(rest)
	}
	if isBarePath(value) {
		return maskLocal(value)
	}
	if isGitValue(value) {
		return maskGit(value)
	}
	if isHTTPURL(value) {
		return MaskedURL
	}

	return MaskedUnknown
}

// MatchesUnsupportedProtocol reports whether the backend treats a resolved value as an unsupported
// protocol, which makes it blank the package version. Those values must stay, in masked form, for
// results to stay the same. All other resolved values aren't used and are removed.
func MatchesUnsupportedProtocol(value string) bool {
	if isBarePath(value) {
		return false
	}
	masked := MaskValue(value)
	switch masked {
	case MaskedGitWord, MaskedGitWordNoUser, MaskedLocalDot, MaskedPatchFile:
		return true
	}

	return strings.HasPrefix(masked, patchProtocol) && strings.HasSuffix(masked, "#"+maskedPatchFileName)
}

func isPublicRegistryURL(value string) bool {
	if !strings.Contains(value, "://") {
		value = "https://" + value
	}
	parsed, err := url.Parse(value)
	if err != nil || (parsed.Scheme != "https" && parsed.Scheme != "http") {
		return false
	}

	// A registry URL with credentials in it is private, whichever host it points to.
	return parsed.User == nil && publicRegistryHosts[parsed.Host]
}

func maskWorkspace(value string) string {
	switch value {
	case "workspace:*", "workspace:^", "workspace:~", "workspace:.":
		return value
	}

	return MaskedWorkspace
}

// maskPatch keeps the descriptor of the patched package, `patch:<name>@<range>#<patch file>`, and
// masks only the patch file. The backend reads the patched package from that descriptor, so
// without it the package and its dependencies change in the scan result. The range inside is
// masked by the same rules as any other value, so a private source can't stay just because a
// patch wraps it.
func maskPatch(value string) string {
	if strings.Contains(value, yarnBuiltinPatchToken) {
		// Built-in patches hold no customer data, apart from a locator naming the project.
		return patchLocatorRegex.ReplaceAllString(value, "")
	}

	descriptor, patchFile, found := strings.Cut(strings.TrimPrefix(value, patchProtocol), "#")
	maskedDescriptor, valid := maskPatchedDescriptor(descriptor)
	switch {
	case found && valid && strings.Contains(patchFile, ".patch"):
		return patchProtocol + maskedDescriptor + "#" + maskedPatchFileName
	case found && valid:
		return patchProtocol + maskedDescriptor + "#" + maskedPatchSource
	case strings.Contains(value, ".patch"):
		return MaskedPatchFile
	}

	return MaskedPatch
}

// maskPatchedDescriptor takes the URL-encoded `<name>@<range>` of a patched package. A range that
// masking doesn't change keeps its original encoding.
func maskPatchedDescriptor(descriptor string) (string, bool) {
	decoded, err := url.PathUnescape(descriptor)
	if err != nil || len(decoded) < 2 {
		return "", false
	}
	index := strings.Index(decoded[1:], "@")
	if index < 0 {
		return "", false
	}
	name, versionRange := decoded[:index+1], decoded[index+2:]
	if !packageNameRegex.MatchString(name) {
		return "", false
	}

	masked := MaskValue(versionRange)
	if masked == versionRange {
		return descriptor, true
	}

	return name + "@" + strings.NewReplacer(":", "%3A", "#", "%23").Replace(masked), true
}

func CutLocalProtocol(value string) (string, bool) {
	for _, protocol := range localProtocols {
		if rest, found := strings.CutPrefix(value, protocol); found {
			return rest, true
		}
	}

	return "", false
}

func isBarePath(value string) bool {
	return strings.HasPrefix(value, ".") || strings.HasPrefix(value, "/") || strings.HasPrefix(value, "~/")
}

func maskLocal(path string) string {
	if strings.HasPrefix(path, ".") {
		return MaskedLocalDot
	}

	return MaskedLocal
}

func isHTTPURL(value string) bool {
	return strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "http://")
}

func isGitValue(value string) bool {
	for _, prefix := range gitPrefixes {
		if strings.HasPrefix(value, prefix) {
			return true
		}
	}
	if isHTTPURL(value) {
		return strings.HasSuffix(stripURLSuffix(value), ".git")
	}

	return gitShorthandRegex.MatchString(value)
}

func maskGit(value string) string {
	wordOnly := wordOnlyRegex.MatchString(gitRepoName(value))
	hasGitUser := strings.Contains(value, "git@")
	switch {
	case wordOnly && hasGitUser:
		return MaskedGitWord
	case wordOnly:
		return MaskedGitWordNoUser
	case hasGitUser:
		return MaskedGitOther
	}

	return MaskedGitOtherNoUser
}

func gitRepoName(value string) string {
	path := stripURLSuffix(value)
	if index := strings.LastIndexAny(path, "/:"); index >= 0 {
		path = path[index+1:]
	}

	return strings.TrimSuffix(path, ".git")
}

func stripURLSuffix(value string) string {
	if index := strings.IndexAny(value, "#?"); index >= 0 {
		return value[:index]
	}

	return value
}
