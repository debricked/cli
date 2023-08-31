package java

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	conf "github.com/debricked/cli/internal/callgraph/config"
	"github.com/debricked/cli/internal/callgraph/finder"
	"github.com/debricked/cli/internal/callgraph/job"
	"github.com/debricked/cli/internal/io"
	"github.com/debricked/cli/internal/tui"
	"github.com/fatih/color"
)

type Strategy struct {
	config     conf.IConfig
	cmdFactory ICmdFactory
	files      []string
	paths      []string
	exclusions []string
	finder     finder.IFinder
	ctx        cgexec.IContext
}

func (s Strategy) Invoke() ([]job.IJob, error) {
	var jobs []job.IJob
	// Filter relevant files

	if s.config == nil {
		strategyWarning("No config is setup")

		return jobs, nil
	}

	var roots []string
	var err error
	// NOTE: Removed to meet cyclic complexity limit of 10.
	// pmConfig := s.config.PackageManager()
	// switch pmConfig {
	// case maven:
	// 	roots, err = s.finder.FindMavenRoots(s.files)
	// default:
	// 	roots, err = s.finder.FindMavenRoots(s.files)
	// }

	roots, err = s.finder.FindMavenRoots(s.files)
	if err != nil {
		strategyWarning("Error while finding roots: " + err.Error())

		return jobs, err
	}

	files := s.files

	if s.config.Build() {
		err = buildProjects(s, roots)
		if err != nil {

			return jobs, err
		}

		// If build, then we need to find the newly built files
		files, _ = s.finder.FindFiles(s.paths, s.exclusions)
	}

	javaClassDirs, _ := s.finder.FindJavaClassDirs(files, false)
	absRoots, _ := finder.ConvertPathsToAbsPaths(roots)
	absClassDirs, _ := finder.ConvertPathsToAbsPaths(javaClassDirs)
	rootClassMapping := finder.MapFilesToDir(absRoots, absClassDirs)

	foundRootsWoClasses := 0
	for _, root := range absRoots {
		if _, ok := rootClassMapping[root]; !ok {
			foundRootsWoClasses += 1
		}
	}
	if foundRootsWoClasses > 0 {
		strategyWarning("Found " + fmt.Sprint(foundRootsWoClasses) + " roots without related classes, make sure to build your project before running.")
	}
	for rootFile, classDirs := range rootClassMapping {
		// For each class paths dir within the root, find GCDPath as entrypoint
		// classDir := finder.GCDPath(classDirs)
		rootDir := filepath.Dir(rootFile)
		jobs = append(jobs, NewJob(
			rootDir,
			classDirs,
			s.cmdFactory,
			io.FileWriter{},
			io.NewArchive(rootDir),
			s.config,
			s.ctx,
			io.FileSystem{},
		),
		)
	}

	return jobs, nil
}

func NewStrategy(config conf.IConfig, files []string, paths []string, exclusions []string, finder finder.IFinder, ctx cgexec.IContext) Strategy {
	return Strategy{config, CmdFactory{}, files, paths, exclusions, finder, ctx}
}

func strategyWarning(errMsg string) {
	err := fmt.Errorf(errMsg)
	warningColor := color.New(color.FgYellow, color.Bold).SprintFunc()
	defaultOutputWriter := log.Writer()
	log.Println(warningColor("Warning: ") + err.Error())
	log.SetOutput(defaultOutputWriter)
}

func buildProjects(s Strategy, roots []string) error {
	spinnerType := "building maven project"
	spinnerManager := tui.NewSpinnerManager("Callgraph Build Project", spinnerType)
	spinnerManager.Start()
	success := false || len(roots) == 0
	errors := []string{}
	for _, rootFile := range roots {
		rootDir := filepath.Dir(rootFile)
		spinner := spinnerManager.AddSpinner(rootDir)
		osCmd, err := s.cmdFactory.MakeBuildMavenCmd(rootDir, s.ctx)
		if err != nil {
			err := "Error while building roots (Make command): " + err.Error() + "\nRoot: " + rootDir
			errors = append(errors, err)
			spinner.Error()
			spinnerManager.SetSpinnerMessage(spinner, rootDir, "fail")

			continue
		}
		cmd := cgexec.NewCommand(osCmd)
		err = cgexec.RunCommand(*cmd, s.ctx)

		if err != nil {
			err := "Error while building roots (Make command): " + err.Error() + "\nRoot: " + rootDir
			errors = append(errors, err)
			spinner.Error()
			spinnerManager.SetSpinnerMessage(spinner, rootDir, "fail")

			continue
		}
		spinnerManager.SetSpinnerMessage(spinner, rootDir, "success")
		spinner.Complete()
		success = true
	}
	spinnerManager.Stop()

	if success {
		return nil
	} else {
		for _, err := range errors {
			strategyWarning(err)
		}

		return fmt.Errorf("Build failed for all projects, if already built disable the build flag")
	}

}
