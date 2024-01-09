package resolution

import (
	"errors"
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/debricked/cli/internal/cmd/cmderror"
	"github.com/debricked/cli/internal/file"
	resolutionFile "github.com/debricked/cli/internal/resolution/file"
	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/strategy"
	"github.com/debricked/cli/internal/tui"
)

var (
	ErrBadOpts = errors.New("failed to type case IOptions")
)

type StrictnessLevel int

const (
	NoFail StrictnessLevel = iota
	FailIfAllFail
	FailIfAnyFail
	FailOrWarn
)

func GetStrictnessLevel(level int) (StrictnessLevel, error) {
	switch level {
	case 0:
		return NoFail, nil
	case 1:
		return FailIfAllFail, nil
	case 2:
		return FailIfAnyFail, nil
	case 3:
		return FailOrWarn, nil
	default:
		return NoFail, fmt.Errorf("invalid strictness level: %d", level)
	}
}

type IResolver interface {
	Resolve(paths []string, options IOptions) (IResolution, error)
}

type Resolver struct {
	finder          file.IFinder
	batchFactory    resolutionFile.IBatchFactory
	strategyFactory strategy.IFactory
	scheduler       IScheduler
	npmPreferred    bool
}

type IOptions interface{}

type DebrickedOptions struct {
	Path                 string
	Exclusions           []string
	Verbose              bool
	Regenerate           int
	NpmPreferred         bool
	ResolutionStrictness StrictnessLevel
}

func NewResolver(
	finder file.IFinder,
	batchFactory resolutionFile.IBatchFactory,
	strategyFactory strategy.IFactory,
	scheduler IScheduler,
) Resolver {
	return Resolver{
		finder,
		batchFactory,
		strategyFactory,
		scheduler,
		false,
	}
}

func (r Resolver) setNpmPreferred(npmPreferred bool) {
	r.batchFactory.SetNpmPreferred(npmPreferred)
}

func (r Resolver) GetExitCode(resolution IResolution, options IOptions) (int, error) {
	dOptions, ok := options.(DebrickedOptions)
	if !ok {
		return 1, ErrBadOpts
	}
	errorCount := resolution.GetJobErrorCount()
	jobCount := len(resolution.Jobs())

	return r.getExitCodeBasedOnStrictness(dOptions.ResolutionStrictness, errorCount, jobCount)
}

func (r Resolver) getExitCodeBasedOnStrictness(strictness StrictnessLevel, errorCount, jobCount int) (int, error) {
	switch strictness {
	case NoFail:
		return r.noFailLogic(errorCount, jobCount)
	case FailIfAllFail:
		return r.failIfAllFailLogic(errorCount, jobCount)
	case FailIfAnyFail:
		return r.failIfAnyFailLogic(errorCount, jobCount)
	case FailOrWarn:
		return r.failOrWarnLogic(errorCount, jobCount)
	default:
		return 0, fmt.Errorf("invalid strictness level: %d", strictness)
	}
}

func (r Resolver) noFailLogic(errorCount, jobCount int) (int, error) {
	return 0, nil
}

func (r Resolver) failIfAllFailLogic(errorCount, jobCount int) (int, error) {
	if errorCount == jobCount {
		return 1, nil
	}

	return 0, nil
}

func (r Resolver) failIfAnyFailLogic(errorCount, jobCount int) (int, error) {
	if errorCount > 0 {
		return 1, nil
	}

	return 0, nil
}

func (r Resolver) failOrWarnLogic(errorCount, jobCount int) (int, error) {
	if errorCount == 0 {
		return 0, nil
	} else if errorCount == jobCount {
		return 1, nil
	}

	return 3, nil
}

func (r Resolver) Resolve(paths []string, options IOptions) (IResolution, error) {
	dOptions, ok := options.(DebrickedOptions)
	if !ok {
		return nil, ErrBadOpts
	}
	files, err := r.refinePaths(paths, dOptions.Exclusions, dOptions.Regenerate)
	if err != nil {
		return nil, err
	}
	r.setNpmPreferred(dOptions.NpmPreferred)
	pmBatches := r.batchFactory.Make(files)

	var jobs []job.IJob
	for _, pmBatch := range pmBatches {
		s, strategyErr := r.strategyFactory.Make(pmBatch, paths)
		if strategyErr == nil {
			newJobs, err := s.Invoke()
			if err != nil {
				return nil, err
			}
			jobs = append(jobs, newJobs...)
		}
	}

	resolution, err := r.scheduler.Schedule(jobs)

	if resolution.HasErr() {
		jobErrList := tui.NewJobsErrorList(os.Stdout, resolution.Jobs())
		renderErr := jobErrList.Render(dOptions.Verbose)
		if renderErr != nil {
			return resolution, renderErr
		}
		code, err := r.GetExitCode(resolution, dOptions)
		if err != nil {
			return resolution, err
		}

		if code != 0 {
			err = cmderror.CommandError{
				Code: code,
				Err:  fmt.Errorf("resolution failed"),
			}

			return resolution, err
		}
	}

	return resolution, err
}

func (r Resolver) refinePaths(paths []string, exclusions []string, regenerate int) ([]string, error) {
	var fileSet = map[string]bool{}
	var dirs []string
	for _, arg := range paths {
		cleanArg := path.Clean(arg)
		if cleanArg == "." {
			dirs = append(dirs, cleanArg)

			continue
		}

		fileInfo, err := os.Stat(arg)
		if err != nil {
			return nil, err
		}

		if fileInfo.IsDir() {
			dirs = append(dirs, path.Clean(arg))
		} else {
			fileSet[path.Clean(arg)] = true
		}
	}

	err := r.searchDirs(fileSet, dirs, exclusions, regenerate)
	if err != nil {
		return nil, err
	}

	var files []string
	for f := range fileSet {
		files = append(files, f)
	}

	return files, nil
}

func (r Resolver) searchDirs(fileSet map[string]bool, dirs []string, exclusions []string, regenerate int) error {
	for _, dir := range dirs {
		err := r.processDir(fileSet, dir, exclusions, regenerate)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r Resolver) processDir(fileSet map[string]bool, dir string, exclusions []string, regenerate int) error {
	fileGroups, err := r.finder.GetGroups(
		dir,
		exclusions,
		false,
		file.StrictAll,
	)
	if err != nil {
		return err
	}
	r.processFileGroups(fileSet, fileGroups, regenerate)

	return nil
}

func (r Resolver) processFileGroups(fileSet map[string]bool, fileGroups file.Groups, regenerate int) {
	for _, fileGroup := range fileGroups.ToSlice() {
		if shouldGenerateLock(fileGroup, regenerate) {
			fileSet[fileGroup.ManifestFile] = true
		}
	}
}

func shouldGenerateLock(fileGroup file.Group, regenerate int) bool {
	if !fileGroup.HasFile() {
		return false
	}
	switch regenerate {
	case 0:
		return !fileGroup.HasLockFiles()
	case 1:
		return onlyNonNativeLockFiles(fileGroup.LockFiles)
	case 2:
		return true
	}

	return false
}

func onlyNonNativeLockFiles(lockFiles []string) bool {
	debrickedLockFilePattern := regexp.MustCompile(`.*\.debricked\.lock`)
	for _, lockFile := range lockFiles {
		if !debrickedLockFilePattern.MatchString(lockFile) {
			return false
		}
	}

	return true

}
