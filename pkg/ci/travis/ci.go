package travis

import (
	"debricked/pkg/ci/env"
	"debricked/pkg/ci/util"
)

const (
	EnvKey      = "TRAVIS"
	Integration = "travis"
)

type Ci struct{}

func (_ Ci) Identify() bool {
	return util.EnvKeyIsSet(EnvKey)
}

func (_ Ci) Parse() (env.Env, error) {
	return env.Env{}, nil
}
