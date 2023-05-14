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
	cmd, _ := CmdFactory{}.MakeMvnCopyDependenciesCmd(workingDir, targetDir)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "mvn")
	assert.Contains(t, args, "-q")
	assert.Contains(t, args, "-B")
	assert.Contains(t, args, "dependency:copy-dependencies")
	assert.Contains(t, args, "-DoutputDirectory=target")
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
	assert.Contains(t, args, "java")
	assert.Contains(t, args, "-jar")
	assert.Contains(t, args, "jarpath")
	assert.Contains(t, args, "targetclasses")
	assert.Contains(t, args, "dependencypath")
}
