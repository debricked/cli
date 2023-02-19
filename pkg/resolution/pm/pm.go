package pm

import (
	"github.com/debricked/cli/pkg/resolution/pm/gomod"
	"github.com/debricked/cli/pkg/resolution/pm/gradle"
	"github.com/debricked/cli/pkg/resolution/pm/maven"
	"github.com/debricked/cli/pkg/resolution/pm/pip"
)

type IPm interface {
	Name() string
	Manifests() []string
}

func Pms() []IPm {
	return []IPm{
		maven.NewPm(),
		gradle.NewPm(),
		gomod.NewPm(),
		pip.NewPm(),
	}
}