package java

import (
	"testing"

	ctxTestdata "github.com/debricked/cli/internal/callgraph/cgexec/testdata"
	"github.com/stretchr/testify/assert"
)

func TestMakeMvnCopyDependenciesCmd(t *testing.T) {
	targetDir := "target"
	ctx, _ := ctxTestdata.NewContextMock()
	cmd, _ := CmdFactory{}.MakeMvnCopyDependenciesCmd(dir, targetDir, ctx)
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
	targetClasses := []string{"targetclasses"}
	dependencyClasses := "dependencypath"
	outputName := ".outputName"
	ctx, _ := ctxTestdata.NewContextMock()
	cmd, err := CmdFactory{}.MakeCallGraphGenerationCmd(jarPath, dir, targetClasses, dependencyClasses, outputName, ctx)

	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "java")
	assert.Contains(t, args, "-jar")
	assert.Contains(t, args, jarPath)
	assert.Contains(t, args, targetClasses[0])
	assert.Contains(t, args, dependencyClasses)
	assert.Contains(t, args, outputName)
}

func TestMakeBuildMavenCmd(t *testing.T) {
	jarPath := "jarpath"
	ctx, _ := ctxTestdata.NewContextMock()
	cmd, err := CmdFactory{}.MakeBuildMavenCmd(jarPath, ctx)

	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "mvn")
	assert.Contains(t, args, "package")
	assert.Contains(t, args, "-q")
	assert.Contains(t, args, "-DskipTests")
}

func TestMakeJavaVersionCmd(t *testing.T) {
	jarPath := "jarpath"
	ctx, _ := ctxTestdata.NewContextMock()
	cmd, err := CmdFactory{}.MakeJavaVersionCmd(jarPath, ctx)

	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "java")
	assert.Contains(t, args, "--version")
}
