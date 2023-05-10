package java

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeGradleCopyDependenciesCmd(t *testing.T) {
	workingDir := "dir"
	gradlew := "gradlew"
	groovyFilePath := "groovyfilename"
	cmd, err := CmdFactory{}.MakeGradleCopyDependenciesCmd(workingDir, gradlew, groovyFilePath)

	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "gradlew")
	assert.Contains(t, args, "groovyfilename")
	assert.ErrorContains(t, err, "executable file not found in")
	assert.ErrorContains(t, err, "PATH")
}

func TestMakeMvnCopyDependenciesCmd(t *testing.T) {
	workingDir := "dir"
	targetDir := "target"
	cmd, err := CmdFactory{}.MakeMvnCopyDependenciesCmd(workingDir, targetDir)

	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "dir")
	assert.Contains(t, args, "target")
}

func TestMakeCallGraphGenerationCmd(t *testing.T) {
	jarPath := "jarpath"
	workingDir := "dir"
	targetClasses := "targetclasses"
	dependencyClasses := "dependencypath"
	cmd, err := CmdFactory{}.MakeCallGraphGenerationCmd(jarPath, workingDir, targetClasses, dependencyClasses)

	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "jarpath")
	assert.Contains(t, args, "dir")
	assert.Contains(t, args, "targetclasses")
	assert.Contains(t, args, "dependencypath")
}
