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
	"github.com/debricked/cli/internal/callgraph/language/java/testdata"
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

	fsMock := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fsMock, zip)

	fs := io.FileSystem{}
	shMock := testdata.MockSootHandler{}

	j := NewJob(dir, files, cmdFactoryMock, writer, archiveMock, config, ctx, fs, shMock)
	assert.Equal(t, []string{"file"}, j.GetFiles())
	assert.Equal(t, "dir", j.GetDir())
	assert.False(t, j.Errors().HasError())
}

func TestRunMakeMavenCopyDependenciesCmdErr(t *testing.T) {
	cmdErr := errors.New("cmd-error")
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	cmdFactoryMock.MvnCopyDepErr = cmdErr
	fileWriterMock := &ioTestData.FileWriterMock{}

	fsMock := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fsMock, zip)

	fs := io.FileSystem{}
	shMock := testdata.MockSootHandler{}

	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven", "")
	ctx, _ := ctxTestdata.NewContextMock()
	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx, fs, shMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), cmdErr)
}

func TestRun(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven", "")
	ctx, _ := ctxTestdata.NewContextMock()

	fsMock := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fsMock, zip)

	fs := io.FileSystem{}
	shMock := testdata.MockSootHandler{}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx, fs, shMock)

	go jobTestdata.WaitStatus(j)
	j.Run()
	fmt.Println("TestRun")
	fmt.Println(j.Errors().GetAll())
	assert.False(t, j.Errors().HasError())
}

func TestRunCallgraphMock(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven", "")
	ctx, _ := ctxTestdata.NewContextMock()
	callgraphMock := testdata.CallgraphMock{}

	fsMock := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fsMock, zip)

	fs := io.FileSystem{}
	shMock := testdata.MockSootHandler{}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx, fs, shMock)
	j.runCallGraph(callgraphMock)

	assert.False(t, j.Errors().HasError())
}

func TestRunCallgraphMockError(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven", "")
	ctx, _ := ctxTestdata.NewContextMock()
	callgraphMock := testdata.CallgraphMock{RunCallGraphWithSetupError: fmt.Errorf("error")}

	fsMock := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fsMock, zip)

	fs := io.FileSystem{}
	shMock := testdata.MockSootHandler{}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx, fs, shMock)
	j.runCallGraph(callgraphMock)

	assert.True(t, j.Errors().HasError())
}

func TestRunPostProcessMock(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven", "")
	ctx, _ := ctxTestdata.NewContextMock()

	fsMock := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs(dir, fsMock, zip)

	fs := io.FileSystem{}
	shMock := testdata.MockSootHandler{}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx, fs, shMock)
	go jobTestdata.WaitStatus(j)
	j.runPostProcess()

	assert.False(t, j.Errors().HasError())
}

func TestRunPostProcessZipFileError(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven", "")
	ctx, _ := ctxTestdata.NewContextMock()
	archiveMock := ioTestData.ArchiveMock{ZipFileError: fmt.Errorf("error")}
	fs := io.FileSystem{}
	shMock := testdata.MockSootHandler{}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx, fs, shMock)
	go jobTestdata.WaitStatus(j)
	j.runPostProcess()

	assert.True(t, j.Errors().HasError())
}

func TestRunPostProcessB64Error(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven", "")
	ctx, _ := ctxTestdata.NewContextMock()
	fs := io.FileSystem{}

	archiveMock := ioTestData.ArchiveMock{B64Error: fmt.Errorf("error")}
	shMock := testdata.MockSootHandler{}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx, fs, shMock)
	go jobTestdata.WaitStatus(j)
	j.runPostProcess()

	assert.True(t, j.Errors().HasError())
}

func TestRunPostProcessCleanupError(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven", "")
	ctx, _ := ctxTestdata.NewContextMock()
	fs := io.FileSystem{}
	shMock := testdata.MockSootHandler{}

	archiveMock := ioTestData.ArchiveMock{CleanupError: fmt.Errorf("error")}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx, fs, shMock)
	go jobTestdata.WaitStatus(j)
	j.runPostProcess()

	assert.True(t, j.Errors().HasError())
}

func TestRunPostProcessCleanupNoFileExistError(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven", "")
	ctx, _ := ctxTestdata.NewContextMock()
	fs := io.FileSystem{}

	err := &os.PathError{}
	err.Err = syscall.ENOENT
	archiveMock := ioTestData.ArchiveMock{CleanupError: err}
	shMock := testdata.MockSootHandler{}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx, fs, shMock)
	go jobTestdata.WaitStatus(j)
	j.runPostProcess()

	assert.False(t, j.Errors().HasError())
}

func TestRunPostProcessFromRoot(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()
	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven", "")
	ctx, _ := ctxTestdata.NewContextMock()
	fs := io.FileSystem{}

	err := &os.PathError{}
	err.Err = syscall.ENOENT
	archiveMock := ioTestData.ArchiveMock{PathError: err, Dir: "."}
	shMock := testdata.MockSootHandler{}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx, fs, shMock)
	go jobTestdata.WaitStatus(j)
	j.runPostProcess()

	jobErrors := j.Errors().GetAll()
	assert.True(t, jobErrors[0] == err)

}

func TestRunWithErrorsIsNotExistFalse(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()

	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven", "")
	ctx, _ := ctxTestdata.NewContextMock()

	fs := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fs, zip)

	fs.IsNotExistBool = false
	shMock := testdata.MockSootHandler{}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx, fs, shMock)
	j.Errors().Critical(fmt.Errorf("error"))

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.True(t, j.Errors().HasError())

	fs.IsNotExistBool = true
}

func TestRunWithErrorsIsNotExistTrue(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	cmdFactoryMock := testdata.NewEchoCmdFactory()

	config := conf.NewConfig("java", nil, map[string]string{"pm": maven}, true, "maven", "")
	ctx, _ := ctxTestdata.NewContextMock()

	fs := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fs, zip)

	fs.IsNotExistBool = true
	shMock := testdata.MockSootHandler{}

	j := NewJob(dir, files, cmdFactoryMock, fileWriterMock, archiveMock, config, ctx, fs, shMock)
	j.Errors().Critical(fmt.Errorf("error"))

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.True(t, j.Errors().HasError())
}
