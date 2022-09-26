package github

import (
	"debricked/pkg/ci/util"
	"os"
	"testing"
)

func TestIdentify(t *testing.T) {
	ci := Ci{}
	value := os.Getenv(EnvKey)
	if util.EnvKeyIsSet(EnvKey) {
		if !ci.Identify() {
			t.Error("failed to assert that CI was identified")
		}
		_ = os.Unsetenv(EnvKey)
		defer os.Setenv(EnvKey, value)

		if ci.Identify() {
			t.Error("failed to assert that CI was not identified")
		}
	} else {
		if ci.Identify() {
			t.Error("failed to assert that CI was not identified")
		}

		_ = os.Setenv(EnvKey, "value")
		defer os.Unsetenv(EnvKey)

		if !ci.Identify() {
			t.Error("failed to assert that CI identified")
		}
	}

}
