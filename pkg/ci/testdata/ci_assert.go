package testdata

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AssertIdentify(t *testing.T, identify func() bool, envKey string) {
	assert.False(t, identify(), "failed to assert that CI was not identified")

	_ = os.Setenv(envKey, "value")
	defer UnsetEnvVar(t, envKey)

	assert.True(t, identify(), "failed to assert that CI was identified")
}
