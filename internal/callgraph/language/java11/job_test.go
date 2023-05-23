package java

import (
	"errors"
	"fmt"
	"os"
	"syscall"
	"testing"

	ctxTestdata "github.com/debricked/cli/internal/callgraph/cgexec/testdata"
	conf "github.com/debricked/cli/internal/callgraph/config"
	jobTestdata "github.com/debricked/cli/internal/callgraph/job/testdata"
	"github.com/debricked/cli/internal/callgraph/language/java11/testdata"
	io "github.com/debricked/cli/internal/io"
	ioTestData "github.com/debricked/cli/internal/io/testdata"
	"github.com/stretchr/testify/assert"
)

const (
	badName = "bad-name"
	dir     = "dir"
)

var files = []string{"file"}

func TestNewJob(t *testing.T) {
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	writer := io.FileWriter{}
	config := conf.Config{}
	ctx, _ := ctxTestdata.NewContextMock()

	fs := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fs, zip)

	j := NewJob(dir, files, cmdFactoryMock, writer, archiveMock, config, ctx)
	assert.Equal(t, []string{"file"}, j.GetFiles())
	assert.Equal(t, "dir", j.GetDir())
	assert.False(t, j.Errors().HasError())
}

func TestRunMakeMavenCopyDependenciesCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MvnCopyDepErr = cmdErr
	fileWriterMock := &ioTestData.FileWriterMock{}

	fs := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fs, zip)

	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven")
	ctx, _ := ctxTestdata.NewContextMock()
	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), cmdErr)
}

func TestRun(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven")
	ctx, _ := ctxTestdata.NewContextMock()

	fs := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fs, zip)

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx)

	go jobTestdata.WaitStatus(j)
	j.Run()

	fmt.Println(j.Errors().GetAll())
	assert.False(t, j.Errors().HasError())
}

func TestRunCallgraphMock(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven")
	ctx, _ := ctxTestdata.NewContextMock()
	callgraphMock := testdata.CallgraphMock{}

	fs := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fs, zip)

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx)
	j.runCallGraph(callgraphMock)

	assert.False(t, j.Errors().HasError())
}

func TestRunCallgraphMockError(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven")
	ctx, _ := ctxTestdata.NewContextMock()
	callgraphMock := testdata.CallgraphMock{RunCallGraphWithSetupError: fmt.Errorf("error")}

	fs := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fs, zip)

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx)
	j.runCallGraph(callgraphMock)

	assert.True(t, j.Errors().HasError())
}

func TestRunPostProcessMock(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven")
	ctx, _ := ctxTestdata.NewContextMock()

	fs := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fs, zip)

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx)
	go jobTestdata.WaitStatus(j)
	j.runPostProcess()

	assert.False(t, j.Errors().HasError())
}

func TestRunPostProcessZipFileError(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven")
	ctx, _ := ctxTestdata.NewContextMock()
	archiveMock := ioTestData.ArchiveMock{ZipFileError: fmt.Errorf("error")}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx)
	go jobTestdata.WaitStatus(j)
	j.runPostProcess()

	assert.True(t, j.Errors().HasError())
}

func TestRunPostProcessB64Error(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven")
	ctx, _ := ctxTestdata.NewContextMock()

	archiveMock := ioTestData.ArchiveMock{B64Error: fmt.Errorf("error")}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx)
	go jobTestdata.WaitStatus(j)
	j.runPostProcess()

	assert.True(t, j.Errors().HasError())
}

func TestRunPostProcessCleanupError(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven")
	ctx, _ := ctxTestdata.NewContextMock()

	archiveMock := ioTestData.ArchiveMock{CleanupError: fmt.Errorf("error")}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx)
	go jobTestdata.WaitStatus(j)
	j.runPostProcess()

	assert.True(t, j.Errors().HasError())
}

func TestRunPostProcessCleanupNoFileExistError(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven")
	ctx, _ := ctxTestdata.NewContextMock()

	err := &os.PathError{}
	err.Err = syscall.ENOENT
	archiveMock := ioTestData.ArchiveMock{CleanupError: err}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx)
	go jobTestdata.WaitStatus(j)
	j.runPostProcess()

	assert.False(t, j.Errors().HasError())
}
