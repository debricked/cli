package golang

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	ctxTestdata "github.com/debricked/cli/internal/callgraph/cgexec/testdata"
	"github.com/debricked/cli/internal/io"
	"github.com/stretchr/testify/assert"
)

func TestCleanSymbol(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "Test with pointer symbol",
			in:   "(*github.com/spf13/afero/mem.File).Open",
			want: "github.com/spf13/afero/mem.File.Open",
		},
		{
			name: "Test with non-pointer symbol",
			in:   "github.com/spf13/afero/mem.File.Open",
			want: "github.com/spf13/afero/mem.File.Open",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanSymbol(tt.in); got != tt.want {
				t.Errorf("cleanSymbol() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCallGraphGeneration(t *testing.T) {
	rootFileDir := filepath.Dir("testdata/fixture/app.go")
	outputName := "debricked-call-graph.golang-test"

	defer func() {
		err := os.Remove("testdata/fixture/debricked-call-graph.golang-test")
		if err != nil {
			fmt.Println(err)
		}
	}()

	ctx, _ := ctxTestdata.NewContextMock()

	tests := []struct {
		name      string
		algorithm string
	}{
		{
			name:      "Test with cha algorithm",
			algorithm: "cha",
		},
		{
			name:      "Test with rta algorithm",
			algorithm: "rta",
		},
		{
			name:      "Test with vta algorithm",
			algorithm: "vta",
		},
		{
			name:      "Test with static algorithm",
			algorithm: "static",
		},
	}

	nodeNames := []string{
		"command-line-arguments.main",
		"command-line-arguments.sayHello",
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cg := NewCallgraphBuilder(
				rootFileDir,
				"app.go",
				outputName,
				io.FileSystem{},
				ctx,
			)
			cg.algorithm = tt.algorithm
			outputPath, err := cg.RunCallGraph()
			assert.Nil(t, err)
			assert.NotEmpty(t, outputPath)
			assert.Equal(t, 2, cg.cgModel.NodeCount())
			for _, nodeName := range nodeNames {
				assert.NotNil(t, cg.cgModel.GetNode(nodeName))
			}
			assert.Equal(t, 1, cg.cgModel.EdgeCount())
			node := cg.cgModel.GetNode("command-line-arguments.sayHello")
			assert.Equal(t, 1, len(node.Parents))

		})
	}
}

func TestIsApplicationNode(t *testing.T) {
	tests := []struct {
		name string
		pwd  string
		in   string
		want bool
	}{
		{
			name: "Test with standard library",
			pwd:  "testdata/fixture",
			in:   "testdata/fixture/main.go.Println",
			want: true,
		},
		{
			name: "Test with non-standard library",
			pwd:  "testdata/fixture",
			in:   "github.com/spf13/afero/mem.File.Open",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsApplicationNode(tt.in, tt.pwd); got != tt.want {
				t.Errorf("IsApplicationNode() = %v, want %v", got, tt.want)
			}
		})
	}

}

func TestRelativeFilename(t *testing.T) {
	tests := []struct {
		name string
		pwd  string
		in   string
		want string
	}{
		{
			name: "Test with standard library",
			pwd:  "testdata/fixture",
			in:   "testdata/fixture/main.go.Println",
			want: "main.go.Println",
		},
		{
			name: "Test with non-standard library",
			pwd:  "testdata/fixture",
			in:   "github.com/spf13/afero/mem.File.Open",
			want: "github.com/spf13/afero/mem.File.Open",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, RelativePath(tt.in, tt.pwd))
		})
	}

}
