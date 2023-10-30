package java

import (
	"fmt"
	"testing"

	"github.com/debricked/cli/internal/callgraph/language/java11/testdata"
	ioTestData "github.com/debricked/cli/internal/io/testdata"
	"github.com/stretchr/testify/assert"
)

func TestRunCallGraphWithSetupMock(t *testing.T) {

	cmdMock := testdata.NewEchoCmdFactory()
	fsMock := ioTestData.FileSystemMock{}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, nil)

	err := cg.RunCallGraphWithSetup()

	assert.Nil(t, err)
}

func TestFsOpenEmbedError(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	fsMock := ioTestData.FileSystemMock{FsOpenEmbedError: fmt.Errorf("error")}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, nil)

	err := cg.RunCallGraphWithSetup()

	assert.NotNil(t, err)
}

func TestMkdirTempError(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	fsMock := ioTestData.FileSystemMock{MkdirTempError: fmt.Errorf("error")}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, nil)

	err := cg.RunCallGraphWithSetup()

	assert.NotNil(t, err)
}

func TestReadAllError(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	fsMock := ioTestData.FileSystemMock{FsReadAllError: fmt.Errorf("error")}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, nil)

	err := cg.RunCallGraphWithSetup()

	assert.NotNil(t, err)
}

func TestWriteFileError(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	fsMock := ioTestData.FileSystemMock{FsWriteFileError: fmt.Errorf("error")}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, nil)

	err := cg.RunCallGraphWithSetup()

	assert.NotNil(t, err)
}

func TestRunCallGraphMock(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	fsMock := ioTestData.FileSystemMock{}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, nil)

	err := cg.RunCallGraph(".")

	assert.Nil(t, err)
}

func TestRunCallGraphErrorMock(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.CallGraphGenErr = fmt.Errorf("error")
	fsMock := ioTestData.FileSystemMock{}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, nil)

	err := cg.RunCallGraph(".")

	assert.NotNil(t, err)
}
