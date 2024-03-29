package testdata

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/debricked/cli/internal/callgraph/job"
	"github.com/debricked/cli/internal/io/err"
	"github.com/stretchr/testify/assert"
)

func AssertPathErr(t *testing.T, jobErrs err.IErrors) {
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
