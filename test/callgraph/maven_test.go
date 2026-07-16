package callgraph

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateCallgraph(t *testing.T) {
	testGenerateCallgraph(t, filepath.Join("testdata", "mvnproj-build"), []string{"callgraph"})
}

func TestGenerateCallgraphSootUp(t *testing.T) {
	testGenerateCallgraph(t, filepath.Join("testdata", "mvnproj-build"), []string{"callgraph", "--java-callgraph-engine", "sootup"})
}

func TestGenerateCallgraphNoBuild(t *testing.T) {
	testGenerateCallgraphNoBuild(t, filepath.Join("testdata", "mvnproj-no-build"), []string{"callgraph", "--no-build"})
}

func TestGenerateCallgraphNoBuildSootUp(t *testing.T) {
	testGenerateCallgraphNoBuild(t, filepath.Join("testdata", "mvnproj-no-build"), []string{"callgraph", "--no-build", "--java-callgraph-engine", "sootup"})
}

func testGenerateCallgraph(t *testing.T, mavenProjectPath string, commandArgs []string) {
	t.Helper()
	requireDebrickedBinary(t)

	tmpFolder := filepath.Join(mavenProjectPath, ".debrickedTmpFolder")
	targetFolder := filepath.Join(mavenProjectPath, "target")
	callgraphFile := filepath.Join(mavenProjectPath, "debricked-call-graph.java")

	assert.NoDirExists(t, tmpFolder)
	assert.NoDirExists(t, targetFolder)
	assert.NoFileExists(t, callgraphFile)

	args := append(commandArgs, mavenProjectPath)
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

func testGenerateCallgraphNoBuild(t *testing.T, mavenProjectPath string, commandArgs []string) {
	t.Helper()
	requireDebrickedBinary(t)

	tmpFolder := filepath.Join(mavenProjectPath, ".debrickedTmpFolder")
	targetFolder := filepath.Join(mavenProjectPath, "target")
	callgraphFile := filepath.Join(mavenProjectPath, "debricked-call-graph.java")

	assert.NoFileExists(t, callgraphFile)
	tmpFolderInfoBefore, tmpErr := os.Stat(tmpFolder)
	tmpFolderExistsBefore := tmpErr == nil
	var tmpFolderModTimeBefore interface{} = nil
	if tmpFolderExistsBefore {
		tmpFolderModTimeBefore = tmpFolderInfoBefore.ModTime()
	} else {
		assert.True(t, os.IsNotExist(tmpErr))
	}

	targetFolderInfoBefore, targetErr := os.Stat(targetFolder)
	require.NoError(t, targetErr)
	targetFolderModTimeBefore := targetFolderInfoBefore.ModTime()

	args := append(commandArgs, mavenProjectPath)
	out, err := exec.Command("debricked", args...).Output()
	fmt.Println("debricked callgraph --no-build output:")
	fmt.Println(string(out))
	assert.NoError(t, err)
	assert.NotContains(t, string(out), "Errors")

	tmpFolderInfoAfter, tmpErrAfter := os.Stat(tmpFolder)
	tmpFolderExistsAfter := tmpErrAfter == nil
	var tmpFolderModTimeAfter interface{} = nil
	if tmpFolderExistsAfter {
		tmpFolderModTimeAfter = tmpFolderInfoAfter.ModTime()
	} else {
		assert.True(t, os.IsNotExist(tmpErrAfter))
	}

	targetFolderInfoAfter, _ := os.Stat(targetFolder)
	targetFolderModTimeAfter := targetFolderInfoAfter.ModTime()
	assert.Equal(t, tmpFolderExistsBefore, tmpFolderExistsAfter)
	assert.True(t, tmpFolderModTimeBefore == tmpFolderModTimeAfter)
	assert.True(t, targetFolderModTimeBefore == targetFolderModTimeAfter)
	assert.FileExists(t, callgraphFile)
	os.Remove(callgraphFile)
}

func requireDebrickedBinary(t *testing.T) {
	t.Helper()

	if _, err := exec.LookPath("debricked"); err != nil {
		t.Skip("debricked binary is not available in PATH")
	}
}
