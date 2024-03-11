package callgraph

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateCallgraph(t *testing.T) {

	mavenProjectPath := filepath.Join("testdata", "mvnproj-build")
	tmpFolder := filepath.Join(mavenProjectPath, ".debrickedTmpFolder")
	targetFolder := filepath.Join(mavenProjectPath, "target")
	callgraphFile := filepath.Join(mavenProjectPath, "debricked-call-graph-java")

	assert.NoDirExists(t, tmpFolder)
	assert.NoDirExists(t, targetFolder)
	assert.NoFileExists(t, callgraphFile)

	args := []string{"callgraph", mavenProjectPath}
	out, err := exec.Command("debricked", args...).Output()
	fmt.Println("debricked callgraph output:")
	fmt.Println(string(out))
	assert.NoError(t, err)
	assert.NotContains(t, string(out), "Errors")

	assert.DirExists(t, tmpFolder)
	assert.DirExists(t, targetFolder)
	assert.FileExists(t, callgraphFile)

	os.RemoveAll(tmpFolder)
	os.RemoveAll(targetFolder)
	os.Remove(callgraphFile)

}

func TestGenerateCallgraphNoBuild(t *testing.T) {

	mavenProjectPath := filepath.Join("testdata", "mvnproj-no-build")
	tmpFolder := filepath.Join(mavenProjectPath, ".debrickedTmpFolder")
	targetFolder := filepath.Join(mavenProjectPath, "target")
	callgraphFile := filepath.Join(mavenProjectPath, "debricked-call-graph-java")

	assert.NoFileExists(t, callgraphFile)
	tmpFolderInfoBefore, tmpErr := os.Stat(tmpFolder)
	assert.NoError(t, tmpErr)
	tmpFolderModTimeBefore := tmpFolderInfoBefore.ModTime()
	targetFolderInfoBefore, targetErr := os.Stat(targetFolder)
	assert.NoError(t, targetErr)
	targetFolderModTimeBefore := targetFolderInfoBefore.ModTime()

	args := []string{"callgraph", mavenProjectPath, "--no-build"}
	out, err := exec.Command("debricked", args...).Output()
	fmt.Println("debricked callgraph --no-build output:")
	fmt.Println(string(out))
	assert.NoError(t, err)
	assert.NotContains(t, string(out), "Errors")

	tmpFolderInfoAfter, _ := os.Stat(tmpFolder)
	tmpFolderModTimeAfter := tmpFolderInfoAfter.ModTime()
	targetFolderInfoAfter, _ := os.Stat(targetFolder)
	targetFolderModTimeAfter := targetFolderInfoAfter.ModTime()
	assert.True(t, tmpFolderModTimeBefore == tmpFolderModTimeAfter)
	assert.True(t, targetFolderModTimeBefore == targetFolderModTimeAfter)
	assert.FileExists(t, callgraphFile)

	os.Remove(callgraphFile)

}
