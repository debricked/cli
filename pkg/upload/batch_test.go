package upload

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/client/testdata"
	"github.com/debricked/cli/pkg/file"
	"github.com/debricked/cli/pkg/git"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestUploadWithBadFiles(t *testing.T) {
	group := file.NewGroup("package.json", nil, []string{"yarn.lock"})
	var groups file.Groups
	groups.Add(*group)
	metaObj, err := git.NewMetaObject("", "repository-name", "commit-name", "", "", "")
	if err != nil {
		t.Fatal("failed to create new MetaObject")
	}

	var c client.IDebClient
	clientMock := testdata.NewDebClientMock()
	mockRes := testdata.MockResponse{
		StatusCode:   http.StatusUnauthorized,
		ResponseBody: nil,
		Error:        errors.New("error"),
	}
	clientMock.AddMockResponse(mockRes)
	clientMock.AddMockResponse(mockRes)
	c = clientMock
	batch := newUploadBatch(&c, groups, metaObj, "CLI")
	output := captureOutput(batch.upload)
	outputAssertions := []string{
		"Failed to upload: package.json",
		"Failed to upload: yarn.lock",
	}
	for _, assertion := range outputAssertions {
		if !strings.Contains(output, assertion) {
			t.Error(fmt.Sprintf("failed to assert that output contained %s", assertion))
		}
	}
}

func TestConcludeWithoutAnyFiles(t *testing.T) {
	batch := newUploadBatch(nil, file.Groups{}, nil, "CLI")
	err := batch.conclude()
	if err == nil {
		t.Error("failed to assert that error occurred")
	}
	if !strings.Contains(err.Error(), "failed to find dependency files") {
		t.Error("failed to asser error message")
	}
}

func captureOutput(f func()) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	f()
	log.SetOutput(os.Stderr)
	return buf.String()
}
