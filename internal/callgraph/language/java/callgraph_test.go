package java

import (
	"fmt"
	"testing"

	"github.com/debricked/cli/internal/callgraph/language/java/testdata"
	ioTestData "github.com/debricked/cli/internal/io/testdata"
	"github.com/stretchr/testify/assert"
)

func TestRunCallGraphWithSetupMock(t *testing.T) {

	cmdMock := testdata.NewEchoCmdFactory()
	fsMock := ioTestData.FileSystemMock{}
	arcMock := ioTestData.ArchiveMock{}
	swMock := testdata.MockSootHandler{}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, arcMock, nil, swMock)

	err := cg.RunCallGraphWithSetup()

	assert.Nil(t, err)
}

func TestRunCallGraphWithSetupSootWrapperError(t *testing.T) {

	cmdMock := testdata.NewEchoCmdFactory()
	fsMock := ioTestData.FileSystemMock{}
	arcMock := ioTestData.ArchiveMock{}
	swMock := testdata.MockSootHandler{GetSootWrapperError: fmt.Errorf("")}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, arcMock, nil, swMock)

	err := cg.RunCallGraphWithSetup()

	assert.Error(t, err)
}

func TestRunCallGraphWithSetupSootVersionError(t *testing.T) {

	cmdMock := testdata.CmdFactoryMock{JavaVersionErr: fmt.Errorf("version error")}
	fsMock := ioTestData.FileSystemMock{}
	arcMock := ioTestData.ArchiveMock{}
	swMock := testdata.MockSootHandler{}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, arcMock, nil, swMock)

	err := cg.RunCallGraphWithSetup()

	assert.Error(t, err)
}

func TestRunCallGraphMock(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	fsMock := ioTestData.FileSystemMock{}
	arcMock := ioTestData.ArchiveMock{}
	swMock := testdata.MockSootHandler{}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, arcMock, nil, swMock)

	err := cg.RunCallGraph(".")

	assert.Nil(t, err)
}

func TestRunCallGraphErrorMock(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.CallGraphGenErr = fmt.Errorf("error")
	fsMock := ioTestData.FileSystemMock{}
	arcMock := ioTestData.ArchiveMock{}
	swMock := testdata.MockSootHandler{}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, arcMock, nil, swMock)

	err := cg.RunCallGraph(".")

	assert.NotNil(t, err)
}
