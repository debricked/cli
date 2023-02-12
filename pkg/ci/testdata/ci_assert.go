package testdata

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func AssertIdentify(t *testing.T, identify func() bool, envKey string) {
	assert.False(t, identify(), "failed to assert that CI was not identified")

	_ = os.Setenv(envKey, "value")
	defer UnsetEnvVar(t, envKey)

	assert.True(t, identify(), "failed to assert that CI was identified")
}
