package upload

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/debricked/cli/pkg/client"
	"github.com/debricked/cli/pkg/client/testdata"
	"github.com/debricked/cli/pkg/file"
	"github.com/debricked/cli/pkg/git"
	"github.com/stretchr/testify/assert"
)

func TestNewDebrickedUploader(t *testing.T) {
	uploader, err := NewUploader(nil)
	assert.Error(t, err)
	assert.Nil(t, uploader)

	var c client.IDebClient = &debClientMock{}
	uploader, err = NewUploader(c)
	assert.NoError(t, err)
	assert.NotNil(t, uploader)
}

func TestUpload(t *testing.T) {
	var c client.IDebClient = &debClientMock{}
	uploader, _ := NewUploader(c)
	metaObject, _ := git.NewMetaObject(
		"testdata/yarn",
		"testdata/yarn",
		"testdata/yarn-commit",
		"",
		"",
		"",
	)
	g := file.NewGroup("testdata/yarn/package.json", nil, []string{"testdata/yarn/yarn.lock"})
	groups := file.Groups{}
	groups.Add(*g)
	uploaderOptions := DebrickedOptions{FileGroups: groups, GitMetaObject: *metaObject, IntegrationsName: "CLI"}
	result, err := uploader.Upload(uploaderOptions)

	assert.NoError(t, err)
	assert.NotNil(t, result)
}

func TestUploadPollingError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	// Create mocked file upload response
	uploadMockRes := testdata.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: io.NopCloser(strings.NewReader("{\"ciUploadId\": 1}")),
	}
	debClientMock.AddMockUriResponse("/api/1.0/open/uploads/dependencies/files", uploadMockRes)

	// Create a mocked finish response
	finishMockRes := testdata.MockResponse{
		StatusCode:   http.StatusNoContent,
		ResponseBody: io.NopCloser(strings.NewReader("{}")),
	}
	debClientMock.AddMockUriResponse("/api/1.0/open/finishes/dependencies/files/uploads", finishMockRes)

	// Create mocked scan result response, 201 is returned when the queue time are too long
	scanMockRes := testdata.MockResponse{
		StatusCode:   http.StatusCreated,
		ResponseBody: io.NopCloser(strings.NewReader("{}")),
	}
	debClientMock.AddMockUriResponse("/api/1.0/open/ci/upload/status", scanMockRes)
	uploader, _ := NewUploader(debClientMock)
	metaObject, _ := git.NewMetaObject(
		"testdata/yarn",
		"testdata/yarn",
		"testdata/yarn-commit",
		"",
		"",
		"",
	)
	g := file.NewGroup("testdata/yarn/package.json", nil, []string{"testdata/yarn/yarn.lock"})
	groups := file.Groups{}
	groups.Add(*g)
	uploaderOptions := DebrickedOptions{FileGroups: groups, GitMetaObject: *metaObject, IntegrationsName: "CLI"}
	result, err := uploader.Upload(uploaderOptions)

	assert.NoError(t, err)
	assert.Nil(t, result)
}

type debClientMock struct{}

func (mock *debClientMock) Post(uri string, _ string, _ *bytes.Buffer) (*http.Response, error) {
	res := &http.Response{
		Status:           "",
		StatusCode:       http.StatusOK,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}
	var resBodyBytes []byte
	if uri == "/api/1.0/open/uploads/dependencies/files" {
		f := uploadedFile{1, 0, 0, 0, "0", 0}
		resBodyBytes, _ = json.Marshal(f)

	} else if uri == "/api/1.0/open/finishes/dependencies/files/uploads" {
		res.StatusCode = http.StatusNoContent
	}

	res.Body = io.NopCloser(strings.NewReader(string(resBodyBytes)))

	return res, nil
}

var progress = 50

func (mock *debClientMock) Get(_ string, _ string) (*http.Response, error) {
	res := &http.Response{
		Status:           "",
		StatusCode:       http.StatusOK,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             nil,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}

	f := &uploadStatus{progress, 0, 0, "", nil, ""}
	progress = progress + progress%100

	resBodyBytes, _ := json.Marshal(f)
	res.Body = io.NopCloser(strings.NewReader(string(resBodyBytes)))

	return res, nil
}
