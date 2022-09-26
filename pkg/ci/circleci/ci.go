package circleci

import (
	"debricked/pkg/ci/env"
	"debricked/pkg/ci/util"
)

const (
	EnvKey      = "CIRCLECI"
	Integration = "circleci"
)

type Ci struct{}

func (_ Ci) Identify() bool {
	return util.EnvKeyIsSet(EnvKey)
}

func (_ Ci) Parse() (env.Env, error) {
	return env.Env{}, nil
}
