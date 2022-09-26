package ci

import "debricked/pkg/ci/env"

type ICi interface {
	Identify() bool
	Map() (env.Env, error)
}
