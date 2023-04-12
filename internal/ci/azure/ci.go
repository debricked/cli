package azure

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/debricked/cli/internal/ci/env"
	"github.com/debricked/cli/internal/ci/util"
)

const (
	EnvKey      = "TF_BUILD"
	Integration = "azureDevOps"
)

type Ci struct{}

func (_ Ci) Identify() bool {
	return util.EnvKeyIsSet(EnvKey)
}

func (_ Ci) Map() (env.Env, error) {
	e := env.Env{}
	owner := filepath.Base(os.Getenv("SYSTEM_COLLECTIONURI"))
	e.Repository = fmt.Sprintf("%s/%s", owner, os.Getenv("BUILD_REPOSITORY_NAME"))
	e.Commit = os.Getenv("BUILD_SOURCEVERSION")
	e.Branch = os.Getenv("BUILD_SOURCEBRANCHNAME")
	e.RepositoryUrl = os.Getenv("BUILD_REPOSITORY_URI")
	e.Integration = Integration
	e.Author = os.Getenv("BUILD_REQUESTEDFOREMAIL")
	e.Filepath = os.Getenv("BUILD_SOURCESDIRECTORY")

	return e, nil
}
