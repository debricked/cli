package strategy

import (
	"fmt"

	"github.com/debricked/cli/pkg/resolution/file"
	"github.com/debricked/cli/pkg/resolution/pm/gomod"
	"github.com/debricked/cli/pkg/resolution/pm/gradle"
	"github.com/debricked/cli/pkg/resolution/pm/maven"
	"github.com/debricked/cli/pkg/resolution/pm/pip"
)

type IFactory interface {
	Make(pmBatch file.IBatch, paths []string) (IStrategy, error)
}

type Factory struct{}

func NewStrategyFactory() Factory {
	return Factory{}
}

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
	default:
		return nil, fmt.Errorf("failed to make strategy from %s", name)
	}
}
