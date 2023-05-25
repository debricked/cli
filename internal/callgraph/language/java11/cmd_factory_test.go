package java

import (
	"os"
	"path"
	"testing"

	"github.com/debricked/cli/internal/callgraph/cgexec"
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
	targetClasses := "targetclasses"
	dependencyClasses := "dependencypath"
	ctx, _ := ctxTestdata.NewContextMock()
	cmd, err := CmdFactory{}.MakeCallGraphGenerationCmd(jarPath, dir, targetClasses, dependencyClasses, ctx)

	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "java")
	assert.Contains(t, args, "-jar")
	assert.Contains(t, args, "jarpath")
	assert.Contains(t, args, "targetclasses")
	assert.Contains(t, args, "dependencypath")
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

func TestMakeBuildMavenCmdFunctional(t *testing.T) {
	workingDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(workingDir)
	javaProjectPath := "testdata/mvnproj"
	javaProjectAbsPath := path.Join(workingDir, javaProjectPath)
	javaProjectTargetAbsPath := path.Join(javaProjectAbsPath, "target")
	assert.NoDirExists(t, javaProjectTargetAbsPath)
	// ctx, _ := ctxTestdata.NewContextMock() // TODO change to real context, no mock
	ctx, _ := cgexec.NewContext(10000)
	cmd, err := CmdFactory{}.MakeBuildMavenCmd(javaProjectAbsPath, ctx)
	if err != nil {
		t.Fatal(err)
	}
	cgexec.RunCommand(cmd, ctx)
	assert.DirExists(t, javaProjectTargetAbsPath)
	os.RemoveAll(javaProjectTargetAbsPath)
	assert.NoDirExists(t, javaProjectTargetAbsPath)
}
