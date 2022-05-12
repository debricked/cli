package tests

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var debrickedVariables = []string{
	"DEBRICKED_SCAN_PATH",
	"DEBRICKED_SCAN_REPOSITORY",
	"DEBRICKED_SCAN_COMMIT",
	"DEBRICKED_SCAN_BRANCH",
	"DEBRICKED_SCAN_AUTHOR",
	"DEBRICKED_SCAN_REPOSITORY_URL",
	"DEBRICKED_SCAN_INTEGRATION",
}

func Test(t *testing.T, ciEnv map[string]string) {
	defer resetEnv(t, ciEnv)
	err := setUpCiEnv(ciEnv)
	if err != nil {
		t.Fatal("failed to set up CI environment variables. Error:", err)
	}

	bytes, err := run()
	output := string(bytes)
	if err != nil {
		t.Error(failureMsg(err))
		t.Fatal(output)
	}
	if len(output) == 0 {
		t.Fatal("failed to assert that output exist")
	}

	err = debrickedEnvIsSetUp(output)
	if err != nil {
		t.Error("failed to assert all Debricked variables were set correctly. Error:", err)
	}
}

func run() ([]byte, error) {
	err := os.Chdir("../scripts")
	if err != nil {
		return nil, err
	}
	return exec.Command("./setup_env.sh").CombinedOutput()
}

func resetEnv(t *testing.T, ciEnv map[string]string) {
	err := os.Chdir("../tests")
	if err != nil {
		t.Error(fmt.Sprintf("failed to reset env. Error: %s", err))
	}

	variables := debrickedVariables
	variables = append(variables, "DEBRICKED_DEBUG")
	for variable, _ := range ciEnv {
		variables = append(variables, variable)
	}

	for _, variable := range variables {
		err = os.Unsetenv(variable)
		if err != nil {
			t.Fatal(fmt.Sprintf("failed to reset env variable: %s. Error: %s", variable, err))
		}
	}
}

func setUpCiEnv(env map[string]string) error {
	env["DEBRICKED_DEBUG"] = "true"
	for variable, value := range env {
		err := os.Setenv(variable, value)
		if err != nil {
			return err
		}
	}

	return nil
}

func debrickedEnvIsSetUp(output string) error {
	var assertions []string
	assertions = append(assertions, debrickedVariables...)
	env := strings.Split(output, "\n\n")[1]
	envVars := strings.Split(env, "\n")

	for _, envVar := range envVars {
		for i := 0; i < len(assertions); i++ {
			variable := assertions[i]
			if strings.Contains(envVar, variable) {
				varPair := strings.Split(envVar, "=")
				if varPair[1] != "" {
					assertions[i] = assertions[len(assertions)-1]
					assertions = assertions[:len(assertions)-1]
				}
			}
		}
	}

	if len(assertions) != 0 {
		return errors.New(
			fmt.Sprintf("missing variables: %s", strings.Join(assertions, ", ")),
		)
	}

	return nil
}

func failureMsg(err error) string {
	return fmt.Sprintf("failed to run setup_env.sh. Error: %s", err)
}
