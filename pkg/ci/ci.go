package ci

import "debricked/pkg/ci/env"

type ICi interface {
	Identify() bool
	Parse() (env.Env, error)
}
