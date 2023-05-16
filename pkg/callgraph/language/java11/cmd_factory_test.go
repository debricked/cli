package java

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeMvnCopyDependenciesCmd(t *testing.T) {
	targetDir := "target"
	cmd, _ := CmdFactory{}.MakeMvnCopyDependenciesCmd(dir, targetDir)
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
	targetClasses := "targetclasses"
	dependencyClasses := "dependencypath"
	cmd, err := CmdFactory{}.MakeCallGraphGenerationCmd(jarPath, dir, targetClasses, dependencyClasses)

	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "java")
	assert.Contains(t, args, "-jar")
	assert.Contains(t, args, "jarpath")
	assert.Contains(t, args, "targetclasses")
	assert.Contains(t, args, "dependencypath")
}
