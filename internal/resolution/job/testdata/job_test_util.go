package testdata

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/stretchr/testify/assert"
)

func AssertPathErr(t *testing.T, jobErrs job.IErrors) {
	var path string
	if runtime.GOOS == "windows" {
		path = "%PATH%"
	} else {
		path = "$PATH"
	}
	errs := jobErrs.GetAll()
	assert.Len(t, errs, 1)
	err := errs[0]
	errMsg := fmt.Sprintf("executable file not found in %s", path)
	assert.ErrorContains(t, err, errMsg)
}

func WaitStatus(j job.IJob) {
	for {
		<-j.ReceiveStatus()
	}
}
