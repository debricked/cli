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
	BadOptsErr = errors.New("failed to type case IOptions")
)

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
	Resolutionstrictness int
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

func (r Resolver) GetExitCode(resolution IResolution, options DebrickedOptions) int {
	errorCount := resolution.GetJobErrorCount()
	jobCount := len(resolution.Jobs())

	switch options.Resolutionstrictness {
	case 0:
		return 0
	case 1:
		if errorCount == jobCount {
			return 1
		}
		return 0
	case 2:
		if errorCount > 0 {
			return 1
		}
		return 0
	case 3:
		if errorCount == 0 {
			return 0
		} else if errorCount == jobCount {
			return 1
		}
		return 3
	default:
		panic(fmt.Sprintf("Invalid strictness level: %d", options.Resolutionstrictness))
	}
}

func (r Resolver) Resolve(paths []string, options IOptions) (IResolution, error) {
	dOptions, ok := options.(DebrickedOptions)
	if !ok {
		return nil, BadOptsErr
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
		code := r.GetExitCode(resolution, dOptions)
		if code != 0 {
			err = cmderror.CommandError{
				Code: code,
				Err:  fmt.Errorf("resolution failed"),
			}
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
		fileGroups, err := r.finder.GetGroups(
			dir,
			exclusions,
			false,
			file.StrictAll,
		)
		if err != nil {
			return err
		}
		for _, fileGroup := range fileGroups.ToSlice() {
			shouldGenerate := shouldGenerateLock(fileGroup, regenerate)
			if shouldGenerate {
				fileSet[fileGroup.ManifestFile] = true
			}
		}
	}

	return nil
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
