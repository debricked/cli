package pm

import (
	"github.com/debricked/cli/internal/resolution/pm/gomod"
	"github.com/debricked/cli/internal/resolution/pm/gradle"
	"github.com/debricked/cli/internal/resolution/pm/maven"
	"github.com/debricked/cli/internal/resolution/pm/pip"
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
