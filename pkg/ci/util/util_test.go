package util

import (
	"os"
	"testing"
)

func TestEnvKeyIsSet(t *testing.T) {
	envKey := "DEBRICKED_CLI_KEY"
	if EnvKeyIsSet(envKey) {
		t.Error("failed to assert that env key was not set")
	}

	_ = os.Setenv(envKey, "")
	if EnvKeyIsSet(envKey) {
		t.Error("failed to assert that env key lacked value")
	}

	_ = os.Setenv(envKey, "value")
	if !EnvKeyIsSet(envKey) {
		t.Error("failed to assert that env key was set")
	}

	err := os.Unsetenv(envKey)
	if err != nil {
		t.Fatal("failed to reset env var")
	}
}
