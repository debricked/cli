package java

import (
	"fmt"
	"testing"

	"github.com/debricked/cli/internal/callgraph/language/java/testdata"
	ioTestData "github.com/debricked/cli/internal/io/testdata"
	"github.com/stretchr/testify/assert"
)

func TestRunCallGraphWithSetupSootWrapperError(t *testing.T) {

	cmdMock := testdata.NewEchoCmdFactory()
	fsMock := ioTestData.FileSystemMock{}
	arcMock := ioTestData.ArchiveMock{}
	shMock := testdata.MockSootHandler{GetSootWrapperError: fmt.Errorf("")}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, arcMock, nil, shMock)

	err := cg.RunCallGraphWithSetup()

	assert.Error(t, err)
}

func TestRunCallGraphWithSetupSootVersionError(t *testing.T) {

	cmdMock := testdata.CmdFactoryMock{JavaVersionErr: fmt.Errorf("version error")}
	fsMock := ioTestData.FileSystemMock{}
	arcMock := ioTestData.ArchiveMock{}
	shMock := testdata.MockSootHandler{}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, arcMock, nil, shMock)

	err := cg.RunCallGraphWithSetup()

	assert.Error(t, err)
}

func TestRunCallGraphMock(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	fsMock := ioTestData.FileSystemMock{}
	arcMock := ioTestData.ArchiveMock{}
	shMock := testdata.MockSootHandler{}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, arcMock, nil, shMock)

	err := cg.RunCallGraph(".")

	assert.Nil(t, err)
}

func TestRunCallGraphErrorMock(t *testing.T) {
	cmdMock := testdata.NewEchoCmdFactory()
	cmdMock.CallGraphGenErr = fmt.Errorf("error")
	fsMock := ioTestData.FileSystemMock{}
	arcMock := ioTestData.ArchiveMock{}
	shMock := testdata.MockSootHandler{}
	cg := NewCallgraph(cmdMock, ".", []string{"."}, ".", ".", fsMock, arcMock, nil, shMock)

	err := cg.RunCallGraph(".")

	assert.NotNil(t, err)
}

func TestIsSootUpWrapperJar(t *testing.T) {
	assert.True(t, isSootUpWrapperJar("/tmp/.debricked/SootUpWrapper.jar"))
	assert.False(t, isSootUpWrapperJar("/tmp/.debricked/soot-wrapper.jar"))
}

func TestExtractSootUpDiagnostics(t *testing.T) {
	stdout := "[SootUpWrapper] Call graph succeeded after excluding 1 problematic jar(s): [foo.jar]\n"
	stderr := "[SootUpWrapper] Warning: TypeAssigner bug in jar 'foo.jar' (x)\n  Excluding jar from deep analysis and retrying...\n"
	lines := extractSootUpDiagnostics(stdout, stderr)

	assert.Len(t, lines, 3)
	assert.Contains(t, lines[0], "Call graph succeeded after excluding")
	assert.Contains(t, lines[1], "TypeAssigner bug in jar")
	assert.Contains(t, lines[2], "Excluding jar from deep analysis and retrying")
}
