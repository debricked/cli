package golang

import (
	"fmt"
	"os"
	"syscall"
	"testing"

	ctxTestdata "github.com/debricked/cli/internal/callgraph/cgexec/testdata"
	conf "github.com/debricked/cli/internal/callgraph/config"
	jobTestdata "github.com/debricked/cli/internal/callgraph/job/testdata"
	"github.com/debricked/cli/internal/callgraph/language/golang/testdata"
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
	writer := io.FileWriter{}
	config := conf.Config{}
	ctx, _ := ctxTestdata.NewContextMock()

	fsMock := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fsMock, zip)

	fs := io.FileSystem{}

	j := NewJob(dir, "main.go", writer, archiveMock, config, ctx, fs)
	assert.Equal(t, []string{"main.go"}, j.GetFiles())
	assert.Equal(t, "dir", j.GetDir())
	assert.False(t, j.Errors().HasError())
}

func TestRun(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	config := conf.NewConfig("golang", nil, nil, true, "go")
	ctx, _ := ctxTestdata.NewContextMock()

	fsMock := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fsMock, zip)

	j := NewJob(dir, "main.go", fileWriterMock, archiveMock, config, ctx, fsMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	fmt.Println(j.Errors().GetAll())
	assert.False(t, j.Errors().HasError())
}

func TestRunCallgraphMockError(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	config := conf.NewConfig("golang", nil, nil, true, "go")
	ctx, _ := ctxTestdata.NewContextMock()
	callgraphMock := testdata.CallgraphMock{RunCallGraphError: fmt.Errorf("error")}

	fsMock := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fsMock, zip)

	fs := io.FileSystem{}

	j := NewJob(dir, "main.go", fileWriterMock, archiveMock, config, ctx, fs)
	j.runCallGraph(callgraphMock)

	assert.True(t, j.Errors().HasError())
}

func TestRunPostProcessZipFileError(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	config := conf.NewConfig("golang", nil, nil, true, "go")
	ctx, _ := ctxTestdata.NewContextMock()
	archiveMock := ioTestData.ArchiveMock{ZipFileError: fmt.Errorf("error")}
	fs := io.FileSystem{}

	j := NewJob(dir, "main.go", fileWriterMock, archiveMock, config, ctx, fs)
	go jobTestdata.WaitStatus(j)
	callgraphMock := testdata.CallgraphMock{}
	j.runCallGraph(callgraphMock)

	assert.True(t, j.Errors().HasError())
}

func TestRunPostProcessB64Error(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	config := conf.NewConfig("golang", nil, nil, true, "go")
	ctx, _ := ctxTestdata.NewContextMock()
	fs := io.FileSystem{}

	archiveMock := ioTestData.ArchiveMock{B64Error: fmt.Errorf("error")}

	j := NewJob(dir, "main.go", fileWriterMock, archiveMock, config, ctx, fs)
	go jobTestdata.WaitStatus(j)
	callgraphMock := testdata.CallgraphMock{}
	j.runCallGraph(callgraphMock)

	assert.True(t, j.Errors().HasError())
}

func TestRunPostProcessCleanupError(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	config := conf.NewConfig("golang", nil, nil, true, "go")
	ctx, _ := ctxTestdata.NewContextMock()
	fs := io.FileSystem{}

	archiveMock := ioTestData.ArchiveMock{CleanupError: fmt.Errorf("error")}

	j := NewJob(dir, "main.go", fileWriterMock, archiveMock, config, ctx, fs)
	go jobTestdata.WaitStatus(j)
	callgraphMock := testdata.CallgraphMock{}
	j.runCallGraph(callgraphMock)

	assert.True(t, j.Errors().HasError())
}

func TestRunPostProcessCleanupNoFileExistError(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	config := conf.NewConfig("golang", nil, nil, true, "go")
	ctx, _ := ctxTestdata.NewContextMock()
	fs := io.FileSystem{}

	err := &os.PathError{}
	err.Err = syscall.ENOENT
	archiveMock := ioTestData.ArchiveMock{CleanupError: err}

	j := NewJob(dir, "main.go", fileWriterMock, archiveMock, config, ctx, fs)
	go jobTestdata.WaitStatus(j)
	callgraphMock := testdata.CallgraphMock{}
	j.runCallGraph(callgraphMock)

	assert.False(t, j.Errors().HasError())
}

func TestRunWithErrorsIsNotExistFalse(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	config := conf.NewConfig("golang", nil, nil, true, "go")
	ctx, _ := ctxTestdata.NewContextMock()

	fs := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fs, zip)

	fs.IsNotExistBool = false

	j := NewJob(dir, "main.go", fileWriterMock, archiveMock, config, ctx, fs)
	j.Errors().Critical(fmt.Errorf("error"))

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.True(t, j.Errors().HasError())

	fs.IsNotExistBool = true
}

func TestRunWithErrorsIsNotExistTrue(t *testing.T) {
	fileWriterMock := &ioTestData.FileWriterMock{}
	config := conf.NewConfig("golang", nil, nil, true, "go")
	ctx, _ := ctxTestdata.NewContextMock()

	fs := ioTestData.FileSystemMock{}
	zip := ioTestData.ZipMock{}
	archiveMock := io.NewArchiveWithStructs("dir", fs, zip)

	fs.IsNotExistBool = true

	j := NewJob(dir, "main.go", fileWriterMock, archiveMock, config, ctx, fs)
	j.Errors().Critical(fmt.Errorf("error"))

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.True(t, j.Errors().HasError())
}
