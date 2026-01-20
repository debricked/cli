package pm

import (
	"github.com/debricked/cli/internal/resolution/pm/bower"
	"github.com/debricked/cli/internal/resolution/pm/composer"
	"github.com/debricked/cli/internal/resolution/pm/gomod"
	"github.com/debricked/cli/internal/resolution/pm/gradle"
	"github.com/debricked/cli/internal/resolution/pm/maven"
	"github.com/debricked/cli/internal/resolution/pm/npm"
	"github.com/debricked/cli/internal/resolution/pm/nuget"
	"github.com/debricked/cli/internal/resolution/pm/pip"
	"github.com/debricked/cli/internal/resolution/pm/poetry"
	"github.com/debricked/cli/internal/resolution/pm/sbt"
	"github.com/debricked/cli/internal/resolution/pm/yarn"
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
		poetry.NewPm(),
		yarn.NewPm(),
		npm.NewPm(),
		bower.NewPm(),
		nuget.NewPm(),
		composer.NewPm(),
		sbt.NewPm(),
	}
}
