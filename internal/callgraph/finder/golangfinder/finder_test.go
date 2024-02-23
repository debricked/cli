package golanfinder

import (
	"testing"

	"github.com/debricked/cli/internal/callgraph/finder"
	"github.com/stretchr/testify/assert"
)

func TestGolangFinderImplementsFinder(t *testing.T) {
	assert.Implements(t, (*finder.IFinder)(nil), new(GolangFinder))
}
