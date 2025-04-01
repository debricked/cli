package strategy

import (
	"fmt"

	"github.com/debricked/cli/internal/resolution/file"
	"github.com/debricked/cli/internal/resolution/pm/bower"
	"github.com/debricked/cli/internal/resolution/pm/composer"
	"github.com/debricked/cli/internal/resolution/pm/gomod"
	"github.com/debricked/cli/internal/resolution/pm/gradle"
	"github.com/debricked/cli/internal/resolution/pm/maven"
	"github.com/debricked/cli/internal/resolution/pm/npm"
	"github.com/debricked/cli/internal/resolution/pm/nuget"
	"github.com/debricked/cli/internal/resolution/pm/pip"
	"github.com/debricked/cli/internal/resolution/pm/sbt"
	"github.com/debricked/cli/internal/resolution/pm/yarn"
)

type IFactory interface {
	Make(pmBatch file.IBatch, paths []string) (IStrategy, error)
}

type Factory struct{}

func NewStrategyFactory() Factory {
	return Factory{}
}

//nolint:all
func (sf Factory) Make(pmFileBatch file.IBatch, paths []string) (IStrategy, error) {
	name := pmFileBatch.Pm().Name()
	switch name {
	case maven.Name:
		return maven.NewStrategy(pmFileBatch.Files()), nil
	case gradle.Name:
		return gradle.NewStrategy(pmFileBatch.Files(), paths), nil
	case gomod.Name:
		return gomod.NewStrategy(pmFileBatch.Files()), nil
	case pip.Name:
		return pip.NewStrategy(pmFileBatch.Files()), nil
	case yarn.Name:
		return yarn.NewStrategy(pmFileBatch.Files()), nil
	case npm.Name:
		return npm.NewStrategy(pmFileBatch.Files()), nil
	case bower.Name:
		return bower.NewStrategy(pmFileBatch.Files()), nil
	case nuget.Name:
		return nuget.NewStrategy(pmFileBatch.Files()), nil
	case composer.Name:
		return composer.NewStrategy(pmFileBatch.Files()), nil
	case sbt.Name:
		return sbt.NewStrategy(pmFileBatch.Files()), nil
	default:
		return nil, fmt.Errorf("failed to make strategy from %s", name)
	}
}
