package gradle

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeFindSubGraphCmd(t *testing.T) {
	cmd, _ := CmdFactory{}.MakeFindSubGraphCmd(".", "gradlew", "init.gradle")
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "gradlew")
	assert.Contains(t, args, "--init-script")
	assert.Contains(t, args, "init.gradle")
	assert.Contains(t, args, "debrickedFindSubProjectPaths")
}

func TestMakeDependenciesGraphCmd(t *testing.T) {
	cmd, _ := CmdFactory{}.MakeDependenciesGraphCmd(".", "gradlew", "init.gradle")
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "gradlew")
	assert.Contains(t, args, "--init-script")
	assert.Contains(t, args, "init.gradle")
	assert.Contains(t, args, "debrickedAllDeps")
}
