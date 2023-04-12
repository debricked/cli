package util

import (
	"os"
	"testing"

	"github.com/debricked/cli/internal/ci/testdata"
	"github.com/stretchr/testify/assert"
)

func TestEnvKeyIsSet(t *testing.T) {
	envKey := "DEBRICKED_CLI_KEY"
	assert.False(t, EnvKeyIsSet(envKey), "failed to assert that env key was not set")

	_ = os.Setenv(envKey, "")
	defer testdata.UnsetEnvVar(t, envKey)

	assert.False(t, EnvKeyIsSet(envKey), "failed to assert that env key lacked value")

	_ = os.Setenv(envKey, "value")
	assert.True(t, EnvKeyIsSet(envKey), "failed to assert that env key was set")
}
