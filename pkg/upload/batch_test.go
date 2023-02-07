package upload

import (
	"bytes"
	"errors"
	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/client/testdata"
	"github.com/debricked/cli/pkg/file"
	"github.com/debricked/cli/pkg/git"
	"io"
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
	var buf bytes.Buffer
	log.SetOutput(&buf)
	err = batch.upload()
	log.SetOutput(os.Stderr)
	output := buf.String()
	if output != "" {
		t.Error("failed to assert that there was no output")
	}
	if err == nil {
		t.Error("failed to assert that an error occurred")
	}
	if !strings.EqualFold("failed to initialize a scan due to badly formatted files", err.Error()) {
		t.Error("failed to assert error message")
	}
}

func TestInitAnalysisWithoutAnyFiles(t *testing.T) {
	batch := newUploadBatch(nil, file.Groups{}, nil, "CLI")
	err := batch.initAnalysis()
	if err == nil {
		t.Error("failed to assert that error occurred")
	}
	if !strings.Contains(err.Error(), "failed to find dependency files") {
		t.Error("failed to asser error message")
	}
}

func TestWaitWithPollingTerminatedError(t *testing.T) {
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
		StatusCode:   http.StatusCreated,
		ResponseBody: io.NopCloser(strings.NewReader("{}")),
	}
	clientMock.AddMockResponse(mockRes)
	c = clientMock
	batch := newUploadBatch(&c, groups, metaObj, "CLI")

	uploadResult, err := batch.wait()

	if uploadResult != nil && err != PollingTerminatedErr {
		t.Fatal("Upload result must be nil and err must be PollingTerminatedErr")
	}
}

func TestInitUploadBadFile(t *testing.T) {
	group := file.NewGroup("testdata/misc/requirements.txt", nil, nil)
	var groups file.Groups
	groups.Add(*group)
	metaObj, err := git.NewMetaObject("", "repository-name", "commit-name", "", "", "")
	if err != nil {
		t.Fatal("failed to create new MetaObject")
	}

	clientMock := testdata.NewDebClientMock()
	mockRes := testdata.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: io.NopCloser(strings.NewReader(`{"message":"An empty file is not allowed"}`)),
	}
	clientMock.AddMockResponse(mockRes)

	var c client.IDebClient = clientMock
	batch := newUploadBatch(&c, groups, metaObj, "CLI")

	files, err := batch.initUpload()
	if len(files) != 0 {
		t.Error("failed to assert that batch could not initialize upload")
	}
	if err == nil {
		t.Error("failed to assert that error occurred")
	}
	if !strings.EqualFold("failed to initialize a scan due to badly formatted files", err.Error()) {
		t.Error("failed to assert error message")
	}
}

func TestInitUpload(t *testing.T) {
	group := file.NewGroup("testdata/yarn/package.json", nil, []string{"testdata/yarn/package.json"})
	var groups file.Groups
	groups.Add(*group)
	metaObj, err := git.NewMetaObject("", "repository-name", "commit-name", "", "", "")
	if err != nil {
		t.Fatal("failed to create new MetaObject")
	}

	clientMock := testdata.NewDebClientMock()
	mockRes := testdata.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: io.NopCloser(strings.NewReader(`{"ciUploadId": 1}`)),
	}
	clientMock.AddMockResponse(mockRes)

	var c client.IDebClient = clientMock
	batch := newUploadBatch(&c, groups, metaObj, "CLI")

	files, err := batch.initUpload()
	if len(files) != 1 {
		t.Error("failed to assert that the init deleted one file from the files to be uploaded")
	}
	if err != nil {
		t.Error("failed to assert that error no occurred")
	}
	if batch.ciUploadId != 1 {
		t.Error("failed to assert ciUploadId")
	}
}
