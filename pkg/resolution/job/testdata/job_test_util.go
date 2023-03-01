package testdata

import (
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/debricked/cli/pkg/resolution/job"
	"github.com/debricked/cli/pkg/resolution/pm/writer"
	"github.com/stretchr/testify/assert"
)

func AssertPathErr(t *testing.T, err error) {
	var path string
	if runtime.GOOS == "windows" {
		path = "%PATH%"
	} else {
		path = "$PATH"
	}
	errMsg := fmt.Sprintf("executable file not found in %s", path)
	assert.ErrorContains(t, err, errMsg)
}

func WaitStatus(j job.IJob) {
	for {
		<-j.ReceiveStatus()
	}
}

func CloseFile(job *job.BaseJob, fileWriter writer.IFileWriter, file *os.File) {
	err := fileWriter.Close(file)
	if err != nil {
		job.Err = err
	}
}
