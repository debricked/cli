// Package sanitize writes redacted copies of dependency files, keeping only what the SCA scan uses.
package sanitize

import (
	"fmt"
	"io/fs"
	"maps"
	"os"
	"path/filepath"

	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/policy"
	"github.com/debricked/cli/internal/sanitize/pm"
	"github.com/debricked/cli/internal/sanitize/pm/util"
)

const (
	// debrickedDir also holds call graph tools and customers' policy files, so only the
	// sanitized folder inside it is ever cleared.
	debrickedDir = ".debricked"
	// OutputDirName is cleared on every run, relative to the scanned path.
	OutputDirName = "sanitized"
	// MapFileName is written outside the output directory, so it can't be uploaded with the files.
	MapFileName = "sanitize-map.json"

	reasonNotSupported = "format is not supported by sanitize"
	reasonNoLockFile   = "manifest has no supported lock file to go with it"
	reasonPolicyFile   = "policy file"
)

// Lock files and manifests are kept apart because a manifest is only useful with a lock file next
// to it: the backend rejects manifests as the dependency file. Every file type that no package
// manager in pm.Pms() lists is skipped and reported, so its original can't reach the output.
var lockFileSanitizers, manifestSanitizers = sanitizersByFileName(pm.Pms())

func sanitizersByFileName(pms []pm.IPm) (map[string]util.FileSanitizer, map[string]util.FileSanitizer) {
	lockFiles := map[string]util.FileSanitizer{}
	manifests := map[string]util.FileSanitizer{}
	for _, packageManager := range pms {
		maps.Copy(lockFiles, packageManager.LockFiles())
		maps.Copy(manifests, packageManager.Manifests())
	}

	return lockFiles, manifests
}

type Options struct {
	RootPath   string
	Exclusions []string
	Inclusions []string
	// LockFileOnly and Strictness filter the files to sanitize, as in `debricked files find`.
	// Their zero values sanitize every file the finder returns.
	LockFileOnly bool
	Strictness   int
}

type FileIssue struct {
	Path   string
	Reason string
}

type Result struct {
	OutputDir string
	MapPath   string
	// Groups has the same shape as the groups from `debricked files find`, with the paths of the
	// sanitized copies. A group is only listed if at least one of its lock files was sanitized.
	Groups    []file.Group
	Sanitized []string
	Skipped   []FileIssue
	Failed    []FileIssue
}

type ISanitizer interface {
	Sanitize(options Options) (Result, error)
}

type Sanitizer struct {
	finder file.IFinder
}

func NewSanitizer(finder file.IFinder) *Sanitizer {
	return &Sanitizer{finder: finder}
}

type run struct {
	root         string
	result       Result
	placeholders util.Placeholders
	seen         map[string]bool
	outputs      map[string]string
	projectCount int
}

// Sanitize writes sanitized copies of the files found under options.RootPath. A file that can't be
// sanitized is listed in the result rather than ending the run, so one bad file doesn't hide the
// others. The returned error is for problems with the run itself.
func (s *Sanitizer) Sanitize(options Options) (Result, error) {
	root := options.RootPath
	if root == "" {
		root = "."
	}
	outputDir := filepath.Join(root, debrickedDir, OutputDirName)
	exclusions := append(append([]string{}, options.Exclusions...), file.SanitizedOutputExclusion)

	groups, err := s.finder.GetGroups(file.DebrickedOptions{
		RootPath:     root,
		Exclusions:   exclusions,
		Inclusions:   options.Inclusions,
		LockFileOnly: options.LockFileOnly,
		Strictness:   options.Strictness,
	})
	if err != nil {
		return Result{}, err
	}

	r := &run{
		root:    root,
		result:  Result{OutputDir: outputDir, MapPath: filepath.Join(root, debrickedDir, MapFileName)},
		seen:    map[string]bool{},
		outputs: map[string]string{},
	}

	// Clearing comes after the search, so a failed search leaves the previous output in place.
	if err = os.RemoveAll(outputDir); err != nil {
		return Result{}, fmt.Errorf("failed to clear %s: %w", outputDir, err)
	}
	if err = os.MkdirAll(outputDir, 0o750); err != nil {
		return Result{}, fmt.Errorf("failed to create %s: %w", outputDir, err)
	}

	for _, group := range groups.ToSlice() {
		r.processGroup(group)
	}
	if err = r.reportPolicyFiles(exclusions, options.Inclusions); err != nil {
		return r.result, err
	}
	if err = util.WriteMap(r.result.MapPath, r.placeholders.Entries()); err != nil {
		return r.result, fmt.Errorf("failed to write the local map: %w", err)
	}

	return r.result, nil
}

func (r *run) processGroup(group file.Group) {
	scope := r.newScopeIfHandled(group)

	sanitizedLockFiles := []string{}
	for _, lockFile := range group.LockFiles {
		if output := r.processFile(lockFile, lockFileSanitizers[filepath.Base(lockFile)], scope); output != "" {
			sanitizedLockFiles = append(sanitizedLockFiles, output)
		}
	}

	sanitizedManifest := ""
	if group.HasFile() {
		sanitizer := manifestSanitizers[filepath.Base(group.ManifestFile)]
		if sanitizer != nil && scope == nil {
			r.skip(group.ManifestFile, reasonNoLockFile)
		} else {
			sanitizedManifest = r.processFile(group.ManifestFile, sanitizer, scope)
		}
	}

	if len(sanitizedLockFiles) > 0 {
		r.result.Groups = append(r.result.Groups, *file.NewGroup(sanitizedManifest, group.CompiledFormat, sanitizedLockFiles))
	}
}

// newScopeIfHandled returns nil for a group with no lock file that sanitize can handle.
func (r *run) newScopeIfHandled(group file.Group) *util.ProjectScope {
	for _, lockFile := range group.LockFiles {
		if lockFileSanitizers[filepath.Base(lockFile)] != nil {
			r.projectCount++

			return util.NewProjectScope(&r.placeholders, r.projectCount)
		}
	}

	return nil
}

// processFile returns the path of the sanitized copy, or an empty string if there is none. It
// handles each path once, because workspace members share their root's lock file.
func (r *run) processFile(path string, sanitizer util.FileSanitizer, scope *util.ProjectScope) string {
	relativePath, firstVisit := r.visit(path)
	if !firstVisit {
		return r.outputs[relativePath]
	}

	if sanitizer == nil {
		r.skip(path, skipReason(path))

		return ""
	}

	content, err := os.ReadFile(path) // #nosec G304 -- path comes from the finder's walk of the scanned directory
	if err != nil {
		r.fail(path, "could not read the file")

		return ""
	}
	sanitized, err := sanitizer(content, util.NewFileContext(scope, relativePath))
	if err != nil {
		r.fail(path, util.ErrInvalidFormat.Error())

		return ""
	}

	destination := filepath.Join(r.result.OutputDir, relativePath)
	if err = os.MkdirAll(filepath.Dir(destination), 0o750); err != nil {
		r.fail(path, "could not create the output folder")

		return ""
	}
	if err = os.WriteFile(destination, sanitized, 0o600); err != nil {
		r.fail(path, "could not write the sanitized copy")

		return ""
	}
	r.result.Sanitized = append(r.result.Sanitized, path)
	r.outputs[relativePath] = destination

	return destination
}

// visit returns the path relative to the scanned directory, and whether this is its first visit.
func (r *run) visit(path string) (string, bool) {
	relativePath, err := filepath.Rel(r.root, path)
	if err != nil {
		relativePath = path
	}
	firstVisit := !r.seen[relativePath]
	r.seen[relativePath] = true

	return relativePath, firstVisit
}

func skipReason(path string) string {
	if filepath.Base(path) == policy.DefaultFileName {
		return reasonPolicyFile
	}

	return reasonNotSupported
}

// reportPolicyFiles walks separately because whether the finder returns a policy file depends on
// the supported formats it gets from the server. Files the finder did return are already reported.
func (r *run) reportPolicyFiles(exclusions, inclusions []string) error {
	return filepath.WalkDir(r.root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() || entry.Name() != policy.DefaultFileName || file.Excluded(exclusions, inclusions, path) {
			return nil
		}
		if _, firstVisit := r.visit(path); firstVisit {
			r.skip(path, reasonPolicyFile)
		}

		return nil
	})
}

func (r *run) skip(path, reason string) {
	r.result.Skipped = append(r.result.Skipped, FileIssue{Path: path, Reason: reason})
}

func (r *run) fail(path, reason string) {
	r.result.Failed = append(r.result.Failed, FileIssue{Path: path, Reason: reason})
}
