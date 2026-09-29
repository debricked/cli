package mcpclient

import (
	"os"
	"runtime"
)

// IEnvironment abstracts OS/environment lookups so client config path resolution is testable
// across platforms without mutating real env vars or touching the real home directory.
type IEnvironment interface {
	Getenv(key string) string
	UserHomeDir() (string, error)
	Getwd() (string, error)
	GOOS() string
	Executable() (string, error)
}

// osEnvironment is the real IEnvironment implementation, backed by the os/runtime packages.
type osEnvironment struct{}

// NewEnvironment creates the real, OS-backed IEnvironment implementation.
func NewEnvironment() IEnvironment {
	return osEnvironment{}
}

func (osEnvironment) Getenv(key string) string {
	return os.Getenv(key)
}

func (osEnvironment) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

func (osEnvironment) Getwd() (string, error) {
	return os.Getwd()
}

func (osEnvironment) GOOS() string {
	return runtime.GOOS
}

func (osEnvironment) Executable() (string, error) {
	return os.Executable()
}
