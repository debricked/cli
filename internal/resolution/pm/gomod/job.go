package gomod

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/debricked/cli/internal/resolution/pm/writer"
)

const (
	fileName                   = "gomod.debricked.lock"
	executableNotFoundErrRegex = `executable file not found`
	versionNotFoundErrRegex    = `require ([^"'\s:]+): version "[^"'\s:]+" invalid: ([^"'\n:]+)`
	revisionNotFoundErrRegex   = `([^"'\s\n:]+): reading [^"'\n:]+ at revision [^"'\n:]+: unknown revision ([^"'\n:]+)`
	dependencyNotFoundErrRegex = `go: ([^"'\s:]+): .*\n.*fatal: could not read Username`
	repositoryNotFoundErrRegex = `go: ([^"'\s:]+): .*\n.*remote: Repository not found`
	noPackageErrRegex          = `([^"'\s:]+): .*, but does not contain package`
	unableToResolveErrRegex    = `go: module ([^"'\s:]+): .*\n.*Permission denied`
	noInternetErrRegex         = `dial tcp: lookup ([^"'\s:]+) .+: server misbehaving`
)

type Job struct {
	job.BaseJob
	cmdFactory ICmdFactory
	fileWriter writer.IFileWriter
}

type PackageDetail struct {
	ImportPath  string   `json:"ImportPath"`
	Imports     []string `json:"Imports"`
	TestImports []string `json:"TestImports"`
}

func NewJob(file string, cmdFactory ICmdFactory, fileWriter writer.IFileWriter) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(file),
		cmdFactory: cmdFactory,
		fileWriter: fileWriter,
	}
}

func (j *Job) Run() {
	status := "creating dependency graph"
	j.SendStatus(status)

	graphCmdOutput, cmd, err := j.runGraphCmd()
	if err != nil {
		j.handleError(j.createError(err.Error(), cmd, status))

		return
	}

	status = "creating dependency version list"
	j.SendStatus(status)
	listCmdOutput, cmd, err := j.runListCmd()
	if err != nil {
		j.handleError(j.createError(err.Error(), cmd, status))

		return
	}

	status = "analyzing package dependencies"
	j.SendStatus(status)
	listJsonOutput, cmd, err := j.runListJsonCmd()
	if err != nil {
		j.handleError(j.createError(err.Error(), cmd, status))

		return
	}

	prodListCmdOutput, devListCmdOutput := j.parseDependencies(listJsonOutput, listCmdOutput)

	status = "creating lock file"
	j.SendStatus(status)
	lockFile, err := j.fileWriter.Create(util.MakePathFromManifestFile(j.GetFile(), fileName))
	if err != nil {
		j.handleError(j.createError(err.Error(), cmd, status))

		return
	}
	defer util.CloseFile(j, j.fileWriter, lockFile)

	var fileContents []byte
	fileContents = append(fileContents, graphCmdOutput...)
	fileContents = append(fileContents, []byte("\n")...)
	fileContents = append(fileContents, prodListCmdOutput...)

	if len(devListCmdOutput) > 0 {
		fileContents = append(fileContents, []byte("\n\n")...)
		fileContents = append(fileContents, devListCmdOutput...)
	}

	err = j.fileWriter.Write(lockFile, fileContents)
	if err != nil {
		j.handleError(j.createError(err.Error(), cmd, status))
	}
}

func (j *Job) getWorkingDir() string {
	return filepath.Dir(filepath.Clean(j.GetFile()))
}

func (j *Job) runGraphCmd() ([]byte, string, error) {
	graphCmd, err := j.cmdFactory.MakeGraphCmd(j.getWorkingDir())
	if err != nil {
		return nil, graphCmd.String(), err
	}

	return j.handleCmdOutput(graphCmd)
}

func (j *Job) runListCmd() ([]byte, string, error) {
	listCmd, err := j.cmdFactory.MakeListCmd(j.getWorkingDir())
	if err != nil {
		return nil, listCmd.String(), err
	}

	return j.handleCmdOutput(listCmd)
}

func (j *Job) runListJsonCmd() ([]byte, string, error) {
	listJsonCmd, err := j.cmdFactory.MakeListJsonCmd(j.getWorkingDir())
	if err != nil {
		return nil, listJsonCmd.String(), err
	}

	return j.handleCmdOutput(listJsonCmd)
}

func (j *Job) handleCmdOutput(cmd *exec.Cmd) ([]byte, string, error) {
	output, err := cmd.Output()
	if err != nil {
		return nil, cmd.String(), j.GetExitError(err, "")
	}

	return output, cmd.String(), nil
}

func (j *Job) getImports(jsonOutput []byte) (map[string]bool, map[string]bool) {
	decoder := json.NewDecoder(bytes.NewReader(jsonOutput))
	imports := make(map[string]bool)
	testImports := make(map[string]bool)

	for {
		var pkg PackageDetail
		if err := decoder.Decode(&pkg); err != nil {
			break
		}

		for _, imported := range pkg.Imports {
			imports[imported] = true
		}

		for _, imported := range pkg.TestImports {
			testImports[imported] = true
		}
	}

	return imports, testImports
}

func (j *Job) parseDependencies(jsonOutput []byte, listCmdOutput []byte) ([]byte, []byte) {
	modules := j.parseModules(listCmdOutput)
	imports, testImports := j.getImports(jsonOutput)

	prodDependencies := make([]string, 0)
	devDependencies := make([]string, 0)
	for dependency, version := range modules {
		depFoundInImports := j.isModuleFound(dependency, imports)
		depFoundInTestImports := j.isModuleFound(dependency, testImports)
		module := strings.TrimSpace(dependency + " " + version)

		if depFoundInTestImports && !depFoundInImports {
			devDependencies = append(devDependencies, module)
		} else {
			prodDependencies = append(prodDependencies, module)
		}
	}

	sort.Strings(prodDependencies)
	sort.Strings(devDependencies)

	return []byte(strings.Join(prodDependencies, "\n")), []byte(strings.Join(devDependencies, "\n"))
}

func (j *Job) isModuleFound(module string, imports map[string]bool) bool {
	result := false

	for importPath := range imports {
		if strings.Contains(importPath, module) {
			result = true

			break
		}
	}

	return result
}

func (j *Job) parseModules(output []byte) map[string]string {
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	modules := make(map[string]string)
	for _, line := range lines {
		parts := strings.Fields(line)
		moduleName := parts[0]
		if len(parts) > 1 {
			version := strings.Join(parts[1:], " ")
			modules[moduleName] = version
		} else {
			modules[moduleName] = ""
		}
	}

	return modules
}

func (j *Job) createError(error string, cmd string, status string) job.IError {
	cmdError := util.NewPMJobError(error)
	cmdError.SetCommand(cmd)
	cmdError.SetStatus(status)

	return cmdError
}

func (j *Job) handleError(cmdError job.IError) {
	expressions := []string{
		executableNotFoundErrRegex,
		versionNotFoundErrRegex,
		revisionNotFoundErrRegex,
		dependencyNotFoundErrRegex,
		repositoryNotFoundErrRegex,
		noPackageErrRegex,
		unableToResolveErrRegex,
		noInternetErrRegex,
	}

	for _, expression := range expressions {
		regex := regexp.MustCompile(expression)
		matches := regex.FindAllStringSubmatch(cmdError.Error(), -1)

		if len(matches) > 0 {
			cmdError = j.addDocumentation(expression, matches, cmdError)
			j.Errors().Append(cmdError)

			return
		}
	}

	j.Errors().Append(cmdError)
}

func (j *Job) addDocumentation(expr string, matches [][]string, cmdError job.IError) job.IError {
	documentation := cmdError.Documentation()

	switch expr {
	case executableNotFoundErrRegex:
		documentation = j.GetExecutableNotFoundErrorDocumentation("Go")
	case versionNotFoundErrRegex:
		documentation = j.getVersionNotFoundErrorDocumentation(matches)
	case revisionNotFoundErrRegex:
		documentation = j.getRevisionNotFoundErrorDocumentation(matches)
	case dependencyNotFoundErrRegex:
		documentation = j.getDependencyNotFoundErrorDocumentation(matches)
	case repositoryNotFoundErrRegex:
		documentation = j.getDependencyNotFoundErrorDocumentation(matches)
	case noPackageErrRegex:
		documentation = j.getNoPackageErrorDocumentation(matches)
	case unableToResolveErrRegex:
		documentation = j.getDependencyNotFoundErrorDocumentation(matches)
	case noInternetErrRegex:
		documentation = j.getNoInternetErrorDocumentation(matches)
	}

	cmdError.SetDocumentation(documentation)

	return cmdError
}

func (j *Job) getVersionNotFoundErrorDocumentation(matches [][]string) string {
	dependency := ""
	recommendation := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		dependency = matches[0][1]
		recommendation = matches[0][2]
	}

	return strings.Join(
		[]string{
			"Failed to find package",
			"\"" + dependency + "\".",
			"Please check that package versions are correct in the manifest file.",
			"It " + recommendation + ".",
		}, " ")
}

func (j *Job) getRevisionNotFoundErrorDocumentation(matches [][]string) string {
	dependency := ""
	revision := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		dependency = matches[0][1]
		revision = matches[0][2]
	}

	return strings.Join(
		[]string{
			"Failed to find package",
			"\"" + dependency + "\" with revision " + revision + ".",
			"Please check that package version is correct in the manifest file.",
		}, " ")
}

func (j *Job) getDependencyNotFoundErrorDocumentation(matches [][]string) string {
	dependency := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		dependency = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Failed to find package",
			"\"" + dependency + "\"",
			"that satisfies the requirements.",
			"Please check that dependencies are correct in the manifest file.",
			"\n" + util.InstallPrivateDependencyMessage,
		}, " ")
}

func (j *Job) getNoPackageErrorDocumentation(matches [][]string) string {
	repository := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		repository = matches[0][1]
	}

	return strings.Join(
		[]string{
			"We weren't able to find a package in provided repository",
			"\"" + repository + "\".",
			"Please check that repository address is spelled correct and it actually contains a Go package.",
		}, " ")
}

func (j *Job) getNoInternetErrorDocumentation(matches [][]string) string {
	registry := ""
	if len(matches) > 0 && len(matches[0]) > 1 {
		registry = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Registry",
			"\"" + registry + "\"",
			"is not available at the moment.",
			"There might be a trouble with your network connection.",
		}, " ")
}
