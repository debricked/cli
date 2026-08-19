package swift

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	executableNotFoundErrRegex = `executable file not found`
	depsFileName               = ".spm.debricked.lock"
)

// dependencyNode mirrors the `swift package show-dependencies --format json`
// tree, which carries every field the backend needs to rebuild the transitive
// dependency tree: package identity, name, source URL, version and children.
type dependencyNode struct {
	Identity     string           `json:"identity"`
	Name         string           `json:"name"`
	URL          string           `json:"url"`
	Version      string           `json:"version"`
	Path         string           `json:"path"`
	Dependencies []dependencyNode `json:"dependencies"`
}

// Job resolves Swift dependencies and emits a debricked lock file containing
// the full dependency tree from `swift package show-dependencies --format json`.
type Job struct {
	job.BaseJob
	cmdFactory ICmdFactory
}

func NewJob(file string, cmdFactory ICmdFactory) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(file),
		cmdFactory: cmdFactory,
	}
}

func (j *Job) Run() {
	status := "generating Package.resolved"
	j.SendStatus(status)

	resolveCmd, err := j.cmdFactory.MakeResolveCmd(j.GetFile())
	if err != nil {
		j.handleError(j.createError(err.Error(), "", status))

		return
	}

	if output, err := resolveCmd.CombinedOutput(); err != nil {
		exitErr := j.GetExitError(err, string(output))
		errorMessage := strings.Join([]string{string(output), exitErr.Error()}, "")
		hint := j.symlinkErrorHint(string(output))
		if hint != "" {
			errorMessage = errorMessage + "\n" + hint
		}
		j.handleError(j.createError(errorMessage, resolveCmd.String(), status))

		return
	}

	status = "generating .spm.debricked.lock"
	j.SendStatus(status)

	depsCmd, err := j.cmdFactory.MakeDepsCmd(j.GetFile())
	if err != nil {
		j.handleError(j.createError(err.Error(), "", status))

		return
	}

	depsOutput, err := depsCmd.CombinedOutput()
	if err != nil {
		exitErr := j.GetExitError(err, string(depsOutput))
		errorMessage := strings.Join([]string{string(depsOutput), exitErr.Error()}, "")
		hint := j.symlinkErrorHint(string(depsOutput))
		if hint != "" {
			errorMessage = errorMessage + "\n" + hint
		}
		j.handleError(j.createError(errorMessage, depsCmd.String(), status))

		return
	}

	status = "writing .spm.debricked.lock"
	j.SendStatus(status)

	tree, err := extractDependencyTree(depsOutput)
	if err != nil {
		j.handleError(j.createError(err.Error(), depsCmd.String(), status))

		return
	}

	err = os.WriteFile(util.MakePathFromManifestFile(j.GetFile(), depsFileName), tree, 0600)
	if err != nil {
		j.handleError(j.createError(err.Error(), "", status))

		return
	}
}

// extractDependencyTree strips non-JSON command noise and verifies that the
// output holds a complete dependency tree before it is persisted for upload.
func extractDependencyTree(output []byte) ([]byte, error) {
	start := bytes.IndexByte(output, '{')
	end := bytes.LastIndexByte(output, '}')
	if start < 0 || end < start {
		return nil, errors.New("swift package show-dependencies did not return a JSON dependency tree")
	}

	tree := bytes.TrimSpace(output[start : end+1])

	var root dependencyNode
	if err := json.Unmarshal(tree, &root); err != nil {
		return nil, err
	}

	if root.Identity == "" && root.Name == "" {
		return nil, errors.New("swift dependency tree is missing root package information")
	}

	return tree, nil
}

func (j *Job) symlinkErrorHint(output string) string {
	if !strings.Contains(output, "unable to create symlink") {
		return ""
	}

	hint := "\nSymlink creation failed. On Windows, this typically requires:\n" +
		"  1. Enable Developer Mode (Settings > Update & Security > For developers)\n" +
		"  2. Run with elevated permissions (Administrator), or\n" +
		"  3. Use WSL2 with Swift toolchain\n" +
		"Reference: https://github.com/apple/swift/issues/61947"

	return hint
}

func (j *Job) createError(errorStr string, cmd string, status string) job.IError {
	cmdError := util.NewPMJobError(errorStr)
	cmdError.SetCommand(cmd)
	cmdError.SetStatus(status)

	return cmdError
}

func (j *Job) handleError(cmdError job.IError) {
	expressions := []string{
		executableNotFoundErrRegex,
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

func (j *Job) addDocumentation(expr string, _ [][]string, cmdError job.IError) job.IError {
	documentation := cmdError.Documentation()

	switch expr {
	case executableNotFoundErrRegex:
		documentation = j.GetExecutableNotFoundErrorDocumentation("Swift")
	}

	cmdError.SetDocumentation(documentation)

	return cmdError
}
