package util

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/job/testdata"
	writerTestdata "github.com/debricked/cli/pkg/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
)

func TestMakePathFromManifestFile(t *testing.T) {
	manifestFile := filepath.Join("pkg", "resolution", "pm", "util", "file.json")
	path := MakePathFromManifestFile(manifestFile, "file.lock")
	lockFile := filepath.Join("pkg", "resolution", "pm", "util", "file.lock")

	assert.Equal(t, lockFile, path)

	path = MakePathFromManifestFile("file.json", "file.lock")
	lockFile = fmt.Sprintf(".%s%s", string(os.PathSeparator), "file.lock")
	assert.Equal(t, lockFile, path)

	path = MakePathFromManifestFile(string(os.PathSeparator), "file.lock")
	assert.Equal(t, fmt.Sprintf("%s%s", string(os.PathSeparator), "file.lock"), path)
}

func TestCloseFile(t *testing.T) {
	var j job.IJob = testdata.NewJobMock("")
	fileWriterMock := writerTestdata.FileWriterMock{}

	CloseFile(j, &fileWriterMock, nil)

	assert.False(t, j.Errors().HasError())
}

func TestCloseFileErr(t *testing.T) {
	var j job.IJob = testdata.NewJobMock("")
	fileWriterMock := writerTestdata.FileWriterMock{}
	closeErr := errors.New("error")
	fileWriterMock.CloseErr = closeErr

	CloseFile(j, &fileWriterMock, nil)

	assert.True(t, j.Errors().HasError())
	criticalErrs := j.Errors().GetCriticalErrors()
	assert.Len(t, criticalErrs, 1)
	criticalErr := criticalErrs[0]
	assert.ErrorIs(t, closeErr, criticalErr)
}
