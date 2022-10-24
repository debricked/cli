package ci

import (
	"errors"
	"github.com/debricked/cli/pkg/ci/argo"
	"github.com/debricked/cli/pkg/ci/azure"
	"github.com/debricked/cli/pkg/ci/bitbucket"
	"github.com/debricked/cli/pkg/ci/buildkite"
	"github.com/debricked/cli/pkg/ci/circleci"
	"github.com/debricked/cli/pkg/ci/env"
	"github.com/debricked/cli/pkg/ci/github"
	"github.com/debricked/cli/pkg/ci/gitlab"
	"github.com/debricked/cli/pkg/ci/travis"
)

type IService interface {
	Find() (env.Env, error)
}

var ErrNotSupported = errors.New("CI is not supported")

type Service struct {
	cis []ICi
}

func NewService(cis []ICi) *Service {
	if cis == nil {
		return &Service{
			[]ICi{
				argo.Ci{},
				azure.Ci{},
				bitbucket.Ci{},
				buildkite.Ci{},
				circleci.Ci{},
				github.Ci{},
				gitlab.Ci{},
				travis.Ci{},
			},
		}
	}
	return &Service{cis}
}

func (s *Service) Find() (env.Env, error) {
	for _, ci := range s.cis {
		if ci.Identify() {
			return ci.Map()
		}
	}

	return env.Env{}, ErrNotSupported
}
