package testdata

import (
	"github.com/debricked/cli/internal/mcpclient"
)

// Environment is a controllable mcpclient.IEnvironment test double.
type Environment struct {
	Env            map[string]string
	HomeDir        string
	HomeDirErr     error
	Wd             string
	WdErr          error
	GOOSValue      string
	ExecutablePath string
	ExecErr        error
}

var _ mcpclient.IEnvironment = Environment{}

func (f Environment) Getenv(key string) string {
	return f.Env[key]
}

func (f Environment) UserHomeDir() (string, error) {
	return f.HomeDir, f.HomeDirErr
}

func (f Environment) Getwd() (string, error) {
	return f.Wd, f.WdErr
}

func (f Environment) GOOS() string {
	return f.GOOSValue
}

func (f Environment) Executable() (string, error) {
	return f.ExecutablePath, f.ExecErr
}
