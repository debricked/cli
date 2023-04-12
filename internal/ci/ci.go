package ci

import "github.com/debricked/cli/internal/ci/env"

type ICi interface {
	Identify() bool
	Map() (env.Env, error)
}
