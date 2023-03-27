package gradle

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGradleSetup(t *testing.T) {

	gs := NewGradleSetup()
	assert.NotNil(t, gs)
}

func TestFindGradleProjectFiles(t *testing.T) {
	gs := NewGradleSetup()
	paths := []string{filepath.Join("testdata", "project")}
	gs.findGradleProjectFiles(paths)

	assert.Len(t, gs.settingsMap, 1)
	assert.Len(t, gs.gradlewMap, 1)
}

func TestErrors(t *testing.T) {

	walkError := GradleSetupWalkError{message: "test"}
	assert.Equal(t, "test", walkError.Error())

	scriptError := GradleSetupScriptError{message: "test"}
	assert.Equal(t, "test", scriptError.Error())

	subprojectError := GradleSetupSubprojectError{message: "test"}
	assert.Equal(t, "test", subprojectError.Error())

}

func TestSetupFilePathMappings(t *testing.T) {
	gs := NewGradleSetup()
	files := []string{filepath.Join("testdata", "project", "build.gradle")}
	gs.setupFilePathMappings(files)

	assert.Len(t, gs.gradlewMap, 1)
	assert.Len(t, gs.settingsMap, 1)
}

func TestSetupFilePathMappingsNoFiles(t *testing.T) {
	gs := NewGradleSetup()
	gs.setupFilePathMappings([]string{})

	assert.Len(t, gs.gradlewMap, 0)
	assert.Len(t, gs.settingsMap, 0)
}

func TestSetupFilePathMappingsNoGradlew(t *testing.T) {
	gs := NewGradleSetup()
	files := []string{filepath.Join("testdata", "project", "subproject", "build.gradle")}
	gs.setupFilePathMappings(files)

	assert.Len(t, gs.gradlewMap, 0)
	assert.Len(t, gs.settingsMap, 0)
}

func TestSetupGradleProjectMappings(t *testing.T) {
	gs := NewGradleSetup()
	gs.settingsMap = map[string]string{
		filepath.Join("testdata", "project", "settings.gradle"): filepath.Join("testdata", "project"),
	}
	gs.subProjectMap = map[string]string{
		filepath.Join("testdata", "project", "subproject", "build.gradle"): filepath.Join("testdata", "project", "subproject"),
	}
	gs.setupGradleProjectMappings()

	assert.Len(t, gs.GradleProjects, 1)
	fmt.Println(gs.GradleProjects)
}

// mock cmd factory
type mockCmdFactory struct {
}

// mock for NewCmd
func (m *mockCmdFactory) MakeFindSubGraphCmd(workingDirectory string, gradlew string, initScript string) (*exec.Cmd, error) {
	//	path, err := exec.LookPath(gradlew)

	// return command that creates a file name ".debricked.multiprojects.txt"
	return exec.Command("ls"), nil
}

func (m *mockCmdFactory) MakeDependenciesCmd(workingDirectory string) (*exec.Cmd, error) {
	return exec.Command("touch", ".debricked.dependencies.txt"), nil
}

func (m *mockCmdFactory) MakeDependenciesGraphCmd(workingDirectory string, gradlew string, initScript string) (*exec.Cmd, error) {
	return &exec.Cmd{
		Path: workingDirectory,
		Args: []string{"touch", ".debricked.dependencies.graph.txt"},
		Dir:  workingDirectory,
	}, nil
}

func TestSetupSubProjectPaths(t *testing.T) {

	gs := NewGradleSetup()
	gs.CmdFactory = &mockCmdFactory{}

	fmt.Println(os.Getwd())

	gradleProject := GradleProject{dir: filepath.Join("testdata", "project"), gradlew: filepath.Join("testdata", "project", "gradlew")}
	gs.setupSubProjectPaths(gradleProject)
	assert.Len(t, gs.subProjectMap, 0)

	gradleProject = GradleProject{dir: filepath.Join("testdata", "project", "subproject"), gradlew: filepath.Join("testdata", "project", "gradlew")}
	gs.setupSubProjectPaths(gradleProject)
	assert.Len(t, gs.subProjectMap, 0)

}

func TestGetGradleW(t *testing.T) {
	gs := NewGradleSetup()
	gs.gradlewMap = map[string]string{
		filepath.Join("testdata", "project"): filepath.Join("testdata", "project", "gradlew"),
	}

	gradlew := gs.GetGradleW(filepath.Join("testdata", "project", "subproject"))

	assert.Equal(t, filepath.Join("testdata", "project", "gradlew"), gradlew)

	gradlew = gs.GetGradleW(filepath.Join("testdata", "project"))

	assert.Equal(t, filepath.Join("testdata", "project", "gradlew"), gradlew)
}
