package sanitize

import (
	"crypto/sha256"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/file/testdata"
	fixtures "github.com/debricked/cli/internal/sanitize/pm/testdata"
	"github.com/debricked/cli/internal/sanitize/pm/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeFixture(t *testing.T, root, relativePath, content string) string {
	t.Helper()
	path := filepath.Join(root, relativePath)
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o750))
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))

	return path
}

func sanitizeGroups(t *testing.T, root string, groups ...*file.Group) Result {
	t.Helper()
	var fileGroups file.Groups
	for _, group := range groups {
		fileGroups.Add(*group)
	}
	finder := testdata.NewFinderMock()
	finder.SetGetGroupsReturnMock(fileGroups, nil)

	result, err := NewSanitizer(finder).Sanitize(Options{RootPath: root})
	require.NoError(t, err)

	return result
}

func readOutput(t *testing.T, result Result, relativePath string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join(result.OutputDir, relativePath)) // #nosec G304 -- reads a file the test run wrote
	require.NoError(t, err)

	return string(content)
}

func readMap(t *testing.T, result Result) []util.MapEntry {
	t.Helper()
	content, err := os.ReadFile(result.MapPath)
	require.NoError(t, err)
	var entries []util.MapEntry
	require.NoError(t, json.Unmarshal(content, &entries))

	return entries
}

func TestSanitizeKeepsRelativePathsAndOriginals(t *testing.T) {
	root := t.TempDir()
	manifest := writeFixture(t, root, "services/portal/package.json", fixtures.PackageJSONBefore)
	lockFile := writeFixture(t, root, "services/portal/yarn.lock", fixtures.YarnV1Before)

	checksum := func(path string) [32]byte {
		content, err := os.ReadFile(path) // #nosec G304 -- reads a fixture the test wrote
		require.NoError(t, err)

		return sha256.Sum256(content)
	}
	manifestBefore, lockFileBefore := checksum(manifest), checksum(lockFile)

	result := sanitizeGroups(t, root, file.NewGroup(manifest, nil, []string{lockFile}))

	assert.Equal(t, manifestBefore, checksum(manifest))
	assert.Equal(t, lockFileBefore, checksum(lockFile))
	assert.ElementsMatch(t, []string{manifest, lockFile}, result.Sanitized)
	assert.Empty(t, result.Failed)
	assert.Equal(t, fixtures.PackageJSONAfter, readOutput(t, result, "services/portal/package.json"))
	assert.Equal(t, fixtures.YarnV1After, readOutput(t, result, "services/portal/yarn.lock"))
}

func TestSanitizeMasksLockFileAndManifestTheSame(t *testing.T) {
	root := t.TempDir()
	source := "git+ssh://git@git.acme.internal/platform/auth-client.git#v2.1.0"
	manifest := writeFixture(t, root, "package.json", `{"dependencies": {"@acme/auth": "`+source+`"}}`)
	lockFile := writeFixture(t, root, "yarn.lock", `"@acme/auth@`+source+`":
  version "2.1.0"
`)

	result := sanitizeGroups(t, root, file.NewGroup(manifest, nil, []string{lockFile}))

	assert.Contains(t, readOutput(t, result, "package.json"), util.MaskedGitOther)
	assert.Contains(t, readOutput(t, result, "yarn.lock"), util.MaskedGitOther)
	assert.NotContains(t, readOutput(t, result, "package.json"), "acme.internal")
	assert.NotContains(t, readOutput(t, result, "yarn.lock"), "acme.internal")
}

func TestSanitizeSkipsAndReportsUnhandledFiles(t *testing.T) {
	root := t.TempDir()
	fingerprints := writeFixture(t, root, "debricked.fingerprints.txt", "s3cret")
	callGraph := writeFixture(t, root, "debricked-call-graph.java", "s3cret")
	pom := writeFixture(t, root, "pom.xml", "<s3cret/>")
	unknownLock := writeFixture(t, root, "unknown-format.lock", "s3cret")

	result := sanitizeGroups(
		t,
		root,
		file.NewGroup("", nil, []string{fingerprints}),
		file.NewGroup(callGraph, nil, []string{}),
		file.NewGroup(pom, nil, []string{}),
		file.NewGroup("", nil, []string{unknownLock}),
	)

	var skipped []string
	for _, issue := range result.Skipped {
		skipped = append(skipped, issue.Path)
	}
	assert.ElementsMatch(t, []string{fingerprints, callGraph, pom, unknownLock}, skipped)
	assert.Empty(t, result.Sanitized)

	entries, err := os.ReadDir(result.OutputDir)
	assert.NoError(t, err)
	assert.Empty(t, entries)
}

const policyFileContent = `{
    "version": "1.0",
    "policies": [
        {
            "name": "No critical vulnerabilities",
            "action": "fail",
            "rules": [
                {
                    "attribute": "cvss3",
                    "operator": ">=",
                    "value": 9
                }
            ]
        }
    ]
}
`

// The finder only returns a policy file when the server lists it as a format, so sanitize also
// finds it on its own. Either way it is copied unchanged and listed once, in a group of its own.
func TestSanitizeCopiesPolicyFiles(t *testing.T) {
	cases := map[string]bool{
		"returned by the finder":     true,
		"not returned by the finder": false,
	}
	for name, returnedByFinder := range cases {
		t.Run(name, func(t *testing.T) {
			root := t.TempDir()
			policyFile := writeFixture(t, root, ".debricked/debricked_policy.json", policyFileContent)
			var groups []*file.Group
			if returnedByFinder {
				groups = append(groups, file.NewGroup(policyFile, nil, []string{}))
			}

			result := sanitizeGroups(t, root, groups...)

			copied := filepath.Join(result.OutputDir, ".debricked", "debricked_policy.json")
			assert.Equal(t, []file.Group{*file.NewGroup(copied, nil, []string{})}, result.Groups)
			assert.Equal(t, []string{policyFile}, result.Sanitized)
			assert.Empty(t, result.Skipped)
			assert.Equal(t, policyFileContent, readOutput(t, result, ".debricked/debricked_policy.json"))
			assert.FileExists(t, policyFile, "files elsewhere in .debricked must survive")
		})
	}
}

func TestSanitizeFailsAPolicyFileThatIsNotJSON(t *testing.T) {
	root := t.TempDir()
	policyFile := writeFixture(t, root, "debricked_policy.json", "version: 1.0\n")

	result := sanitizeGroups(t, root)

	assert.Equal(t, []FileIssue{{Path: policyFile, Reason: util.ErrInvalidFormat.Error()}}, result.Failed)
	assert.Empty(t, result.Groups)
	assert.NoFileExists(t, filepath.Join(result.OutputDir, "debricked_policy.json"))
}

func TestSanitizeSkipsManifestWithoutSupportedLockFile(t *testing.T) {
	root := t.TempDir()
	alone := writeFixture(t, root, "alone/package.json", fixtures.PackageJSONBefore)
	withUnsupported := writeFixture(t, root, "other/package.json", fixtures.PackageJSONBefore)
	unsupportedLock := writeFixture(t, root, "other/unknown-format.lock", "s3cret\n")

	result := sanitizeGroups(
		t,
		root,
		file.NewGroup(alone, nil, []string{}),
		file.NewGroup(withUnsupported, nil, []string{unsupportedLock}),
	)

	assert.Empty(t, result.Sanitized)
	assert.Len(t, result.Skipped, 3)
	assert.NoFileExists(t, filepath.Join(result.OutputDir, "alone/package.json"))
	assert.NoFileExists(t, filepath.Join(result.OutputDir, "other/unknown-format.lock"))
}

func TestSanitizeNeverWritesAFileThatFailed(t *testing.T) {
	root := t.TempDir()
	leakedURL := "https://ci-user:s3cret@artifactory.acme.internal" // #nosec G101 -- fake credentials, the test checks they never reach the output
	broken := writeFixture(t, root, "broken/yarn.lock", leakedURL+"\n")
	good := writeFixture(t, root, "good/yarn.lock", fixtures.YarnV1Before)

	result := sanitizeGroups(
		t,
		root,
		file.NewGroup("", nil, []string{broken}),
		file.NewGroup("", nil, []string{good}),
	)

	require.Len(t, result.Failed, 1)
	assert.Equal(t, broken, result.Failed[0].Path)
	assert.NotContains(t, result.Failed[0].Reason, "s3cret")
	assert.NoFileExists(t, filepath.Join(result.OutputDir, "broken/yarn.lock"))
	assert.FileExists(t, filepath.Join(result.OutputDir, "good/yarn.lock"))
}

func TestSanitizeClearsOnlyTheOutputDirectory(t *testing.T) {
	root := t.TempDir()
	stale := writeFixture(t, root, ".debricked/sanitized/old/yarn.lock", "stale")
	tool := writeFixture(t, root, ".debricked/soot-wrapper.jar", "tool")

	sanitizeGroups(t, root)

	assert.NoFileExists(t, stale)
	assert.FileExists(t, tool)
}

func TestSanitizeLeavesOutputWhenFinderFails(t *testing.T) {
	root := t.TempDir()
	previous := writeFixture(t, root, ".debricked/sanitized/yarn.lock", "previous")
	finder := testdata.NewFinderMock()
	finder.SetGetGroupsReturnMock(file.Groups{}, errors.New("finder failed"))

	_, err := NewSanitizer(finder).Sanitize(Options{RootPath: root})

	assert.Error(t, err)
	assert.FileExists(t, previous)
}

func TestSanitizeHandlesSharedWorkspaceLockFileOnce(t *testing.T) {
	root := t.TempDir()
	rootManifest := writeFixture(t, root, "package.json", fixtures.PackageJSONBefore)
	memberManifest := writeFixture(t, root, "packages/a/package.json", fixtures.PackageJSONBefore)
	lockFile := writeFixture(t, root, "yarn.lock", fixtures.YarnBerryBefore)

	result := sanitizeGroups(
		t,
		root,
		file.NewGroup(rootManifest, nil, []string{lockFile}),
		file.NewGroup(memberManifest, nil, []string{lockFile}),
	)

	assert.Len(t, result.Sanitized, 3)
	assert.Len(t, readMap(t, result), 1)
}

func TestSanitizeReturnsGroupsWithSanitizedPaths(t *testing.T) {
	root := t.TempDir()
	rootManifest := writeFixture(t, root, "package.json", fixtures.PackageJSONBefore)
	memberManifest := writeFixture(t, root, "packages/a/package.json", fixtures.PackageJSONBefore)
	lockFile := writeFixture(t, root, "yarn.lock", fixtures.YarnBerryBefore)
	paket := writeFixture(t, root, "net/paket.lock", fixtures.PaketBefore)
	brokenManifest := writeFixture(t, root, "broken/package.json", fixtures.PackageJSONBefore)
	brokenLock := writeFixture(t, root, "broken/yarn.lock", "not a lock file\n")
	pom := writeFixture(t, root, "java/pom.xml", "<project/>")

	result := sanitizeGroups(
		t,
		root,
		file.NewGroup(rootManifest, nil, []string{lockFile}),
		file.NewGroup(memberManifest, nil, []string{lockFile}),
		file.NewGroup("", nil, []string{paket}),
		file.NewGroup(brokenManifest, nil, []string{brokenLock}),
		file.NewGroup(pom, nil, []string{}),
	)

	output := func(relativePath string) string {
		return filepath.Join(result.OutputDir, relativePath)
	}
	// The member lists the shared lock file, as find does, though only one copy is written.
	assert.Equal(t, []file.Group{
		*file.NewGroup(output("package.json"), nil, []string{output("yarn.lock")}),
		*file.NewGroup(output("packages/a/package.json"), nil, []string{output("yarn.lock")}),
		*file.NewGroup("", nil, []string{output("net/paket.lock")}),
	}, result.Groups)
	for _, group := range result.Groups {
		for _, path := range group.GetAllFiles() {
			assert.FileExists(t, path)
		}
	}
}

func TestSanitizeMapIsOutsideOutputDirectory(t *testing.T) {
	root := t.TempDir()
	lockFile := writeFixture(t, root, "yarn.lock", fixtures.YarnBerryBefore)

	result := sanitizeGroups(t, root, file.NewGroup("", nil, []string{lockFile}))

	relative, err := filepath.Rel(result.OutputDir, result.MapPath)
	assert.NoError(t, err)
	assert.True(t, strings.HasPrefix(relative, ".."), "map at %s", relative)

	assert.Equal(t, []util.MapEntry{{
		Placeholder: "debricked-redacted-project-1",
		Name:        "@acme/payments-portal",
		Path:        "yarn.lock",
	}}, readMap(t, result))

	// Windows has no Unix permission bits: Go reports 0o666 for any file that isn't read-only.
	if runtime.GOOS != "windows" {
		info, err := os.Stat(result.MapPath)
		assert.NoError(t, err)
		assert.Equal(t, os.FileMode(0o600), info.Mode().Perm())
	}

	assert.NotContains(t, readOutput(t, result, "yarn.lock"), "payments-portal")
}

// /sbom/merge matches components by purl, so projects stay separate as long as their placeholders
// differ. This covers the CLI's side of that. The endpoint itself isn't available to this test.
func TestSanitizeKeepsTwoProjectsWithTheSameNameSeparate(t *testing.T) {
	root := t.TempDir()
	project := func(dependency string) string {
		return `"portal@workspace:.":
  version: 0.0.0-use.local
  dependencies:
    ` + dependency + `: "npm:^1.0.0"
`
	}
	first := writeFixture(t, root, "first/yarn.lock", project("express"))
	second := writeFixture(t, root, "second/yarn.lock", project("koa"))

	result := sanitizeGroups(
		t,
		root,
		file.NewGroup("", nil, []string{first}),
		file.NewGroup("", nil, []string{second}),
	)

	firstOutput, secondOutput := readOutput(t, result, "first/yarn.lock"), readOutput(t, result, "second/yarn.lock")
	assert.Contains(t, firstOutput, "debricked-redacted-project-1@workspace:.")
	assert.Contains(t, firstOutput, "express")
	assert.NotContains(t, firstOutput, "koa")
	assert.Contains(t, secondOutput, "debricked-redacted-project-2@workspace:.")
	assert.Contains(t, secondOutput, "koa")
	assert.NotContains(t, secondOutput, "express")

	assert.Equal(t, []util.MapEntry{
		{Placeholder: "debricked-redacted-project-1", Name: "portal", Path: filepath.Join("first", "yarn.lock")},
		{Placeholder: "debricked-redacted-project-2", Name: "portal", Path: filepath.Join("second", "yarn.lock")},
	}, readMap(t, result))
}

func TestSanitizeNumbersPathPlaceholdersAcrossTheRun(t *testing.T) {
	root := t.TempDir()
	content := `{"lockfileVersion": 3, "packages": {"../utils": {"version": "1.0.0"}}}`
	first := writeFixture(t, root, "first/package-lock.json", content)
	second := writeFixture(t, root, "second/package-lock.json", content)

	result := sanitizeGroups(
		t,
		root,
		file.NewGroup("", nil, []string{first}),
		file.NewGroup("", nil, []string{second}),
	)

	assert.Contains(t, readOutput(t, result, "first/package-lock.json"), `"./redacted-1"`)
	assert.Contains(t, readOutput(t, result, "second/package-lock.json"), `"./redacted-2"`)
	assert.Equal(t, []util.MapEntry{
		{Placeholder: "./redacted-1", Name: "../utils", Path: filepath.Join("first", "package-lock.json")},
		{Placeholder: "./redacted-2", Name: "../utils", Path: filepath.Join("second", "package-lock.json")},
	}, readMap(t, result))
}

func TestSanitizePaketFilePlaceholdersUseTheProjectNumber(t *testing.T) {
	root := t.TempDir()
	content := "GITHUB\n  remote: acme/scripts\n    src/Release.fs (9f3c1ab)\n    src/Other.fs (9f3c1ab)\n"
	first := writeFixture(t, root, "first/paket.lock", content)
	second := writeFixture(t, root, "second/paket.lock", content)

	result := sanitizeGroups(
		t,
		root,
		file.NewGroup("", nil, []string{first}),
		file.NewGroup("", nil, []string{second}),
	)

	assert.Equal(t,
		"    debricked-redacted-file-1-1 (9f3c1ab)\n    debricked-redacted-file-1-2 (9f3c1ab)\n",
		readOutput(t, result, "first/paket.lock"),
	)
	assert.Equal(t,
		"    debricked-redacted-file-2-1 (9f3c1ab)\n    debricked-redacted-file-2-2 (9f3c1ab)\n",
		readOutput(t, result, "second/paket.lock"),
	)
	assert.Len(t, readMap(t, result), 4)
}

func TestSanitizeWritesEmptyMapWhenNothingWasRenamed(t *testing.T) {
	result := sanitizeGroups(t, t.TempDir())

	assert.Empty(t, readMap(t, result))
}

type recordingFinder struct {
	*testdata.FinderMock
	options file.DebrickedOptions
}

func (f *recordingFinder) GetGroups(options file.DebrickedOptions) (file.Groups, error) {
	f.options = options

	return f.FinderMock.GetGroups(options)
}

// A custom --exclusion list replaces the default one, so sanitize adds its own output folder itself.
// Otherwise a second run would sanitize the copies from the first.
func TestSanitizeAlwaysExcludesItsOutputDirectory(t *testing.T) {
	cases := map[string][]string{
		"no exclusions":     nil,
		"custom exclusions": {"**/node_modules/**"},
	}
	for name, exclusions := range cases {
		t.Run(name, func(t *testing.T) {
			finder := &recordingFinder{FinderMock: testdata.NewFinderMock()}

			_, err := NewSanitizer(finder).Sanitize(Options{RootPath: t.TempDir(), Exclusions: exclusions})

			require.NoError(t, err)
			assert.Contains(t, finder.options.Exclusions, file.SanitizedOutputExclusion)
			for _, exclusion := range exclusions {
				assert.Contains(t, finder.options.Exclusions, exclusion)
			}
		})
	}
}

func TestSanitizerForMatchesFileNamePatterns(t *testing.T) {
	cases := map[string]bool{
		"Cargo.lock":                 true,
		"sub/maven_install.json":     true,
		"rules_jvm_install.json":     true,
		"gems.Gemfile.lock":          true,
		"Gemfile.lock":               true,
		"install.json":               false,
		"Package.swift":              false,
		"Cargo.toml":                 false,
		"maven_install.json.backup":  false,
		"debricked.fingerprints.txt": false,
	}
	for path, handled := range cases {
		assert.Equal(t, handled, sanitizerFor(lockFileSanitizers, path) != nil, path)
	}
}
