package gradle

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	writerTestdata "github.com/debricked/cli/pkg/resolution/pm/writer/testdata"

	"github.com/debricked/cli/pkg/resolution/pm/writer"
	"github.com/stretchr/testify/assert"
)

func TestNewGradleSetup(t *testing.T) {

	gs := NewGradleSetup()
	assert.NotNil(t, gs)
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
	gs.CmdFactory = &mockCmdFactory{}

	gs.settingsMap = map[string]string{
		filepath.Join("testdata", "project"): filepath.Join("testdata", "project", "settings.gradle"),
	}
	gs.subProjectMap = map[string]string{}
	err := gs.setupGradleProjectMappings()
	// assert GradleSetupSubprojectError
	assert.NotNil(t, err)

	assert.Len(t, gs.GradleProjects, 1)
}

type mockCmdFactory struct {
}

func (m *mockCmdFactory) MakeFindSubGraphCmd(workingDirectory string, gradlew string, initScript string) (*exec.Cmd, error) {
	fileName := filepath.Join(workingDirectory, ".debricked.multiprojects.txt")
	content := []byte(workingDirectory)
	file, err := os.Create(fileName)
	if err != nil {

		return nil, err
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {

		return nil, err
	}

	// if windows use dir
	if runtime.GOOS == "windows" {
		// gradlewOsName = "gradlew.bat"
		exec.Command("dir")
	}

	return exec.Command("ls"), nil
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

	absPath, _ := filepath.Abs(filepath.Join("testdata", "project"))
	gradleProject := GradleProject{dir: absPath, gradlew: filepath.Join("testdata", "project", "gradlew")}
	err := gs.setupSubProjectPaths(gradleProject)
	assert.Nil(t, err)
	assert.Len(t, gs.subProjectMap, 1)

	absPath, _ = filepath.Abs(filepath.Join("testdata", "project", "subproject"))
	gradleProject = GradleProject{dir: absPath, gradlew: filepath.Join("testdata", "project", "gradlew")}
	err = gs.setupSubProjectPaths(gradleProject)
	assert.Nil(t, err)
	assert.Len(t, gs.subProjectMap, 2)

}

func TestSetupSubProjectPathsError(t *testing.T) {

	gs := NewGradleSetup()

	absPath, _ := filepath.Abs(filepath.Join("testdata", "project"))
	gradleProject := GradleProject{dir: absPath, gradlew: filepath.Join("testdata", "project", "gradlew")}
	err := gs.setupSubProjectPaths(gradleProject)

	// assery GradleSetupSubprojectError
	assert.NotNil(t, err)
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

type mockInitFileHandler struct{}

func (_ mockInitFileHandler) ReadInitFile() ([]byte, error) {

	return gradleInitScript.ReadFile("gradle-init/gradle-init-script.groovy")
}

func (i mockInitFileHandler) WriteInitFile(targetFileName string, fileWriter writer.IFileWriter) error {

	return GradleSetupScriptError{message: "read-error"}
}

type mockFileFinder struct{}

func (f mockFileFinder) FindGradleProjectFiles(paths []string) (map[string]string, map[string]string, error) {

	return nil, nil, GradleSetupWalkError{message: "mock error"}
}

// test setup
func TestSetupErrors(t *testing.T) {
	gs := NewGradleSetup()
	gs.Writer = &writerTestdata.FileWriterMock{}
	_, err := gs.Setup([]string{"testdata/project"}, []string{"testdata/project"})
	assert.NotNil(t, err)

	gs.FileFinder = mockFileFinder{}
	_, err = gs.Setup([]string{"testdata/project"}, []string{"testdata/project"})
	assert.Equal(t, "mock error", err.Error())

	gs.InitFileHandler = mockInitFileHandler{}
	_, err = gs.Setup([]string{"testdata/project"}, []string{"testdata/project"})
	assert.Equal(t, "read-error", err.Error())
}
