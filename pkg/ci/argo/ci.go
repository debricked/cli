package argo

import (
	"debricked/pkg/ci/env"
	"debricked/pkg/ci/util"
)

const EnvKey = "BUILDKITE"
const Integration = "argoWorkflows"

type Ci struct{}

func (_ Ci) Identify() bool {
	return util.EnvKeyIsSet(EnvKey)
}

func (_ Ci) Parse() (env.Env, error) {
	return env.Env{}, nil
}
