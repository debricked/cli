package upload

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/debricked/cli/internal/client"
	"github.com/debricked/cli/internal/client/testdata"
	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/git"
	"github.com/stretchr/testify/assert"
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
	batch := newUploadBatch(&c, groups, metaObj, "CLI", 10*60, true, &DebrickedConfig{})
	var buf bytes.Buffer
	log.SetOutput(&buf)
	err = batch.upload()
	log.SetOutput(os.Stderr)
	output := buf.String()

	assert.Empty(t, output)
	assert.ErrorContains(t, err, "Failed to initialize a scan for")
}

func TestInitAnalysisWithoutAnyFiles(t *testing.T) {
	batch := newUploadBatch(nil, file.Groups{}, nil, "CLI", 10*60, true, &DebrickedConfig{})
	err := batch.initAnalysis()

	assert.ErrorContains(t, err, "failed to find dependency files")
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
	batch := newUploadBatch(&c, groups, metaObj, "CLI", 10*60, true, &DebrickedConfig{})

	uploadResult, err := batch.wait()

	assert.Nil(t, uploadResult)
	assert.ErrorIs(t, err, PollingTerminatedErr)
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
	batch := newUploadBatch(&c, groups, metaObj, "CLI", 10*60, true, &DebrickedConfig{})

	files, err := batch.initUpload()

	assert.Empty(t, files)
	assert.ErrorContains(t, err, "Failed to initialize a scan for")
	assert.ErrorContains(t, err, "testdata/misc/requirements.txt")
	assert.ErrorContains(t, err, "tried to upload empty file")
}

func TestInitUploadFingerprintsFree(t *testing.T) {
	group := file.NewGroup("testdata/misc/debricked.fingerprints.txt", nil, nil)
	var groups file.Groups
	groups.Add(*group)
	metaObj, err := git.NewMetaObject("", "repository-name", "commit-name", "", "", "")
	if err != nil {
		t.Fatal("failed to create new MetaObject")
	}

	clientMock := testdata.NewDebClientMock()
	clientMock.SetEnterpriseCustomer(false)
	var c client.IDebClient = clientMock
	batch := newUploadBatch(&c, groups, metaObj, "CLI", 10*60, true, &DebrickedConfig{})

	files, err := batch.initUpload()

	assert.Empty(t, files)
	assert.ErrorContains(t, err, "non-enterprise")
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
	batch := newUploadBatch(&c, groups, metaObj, "CLI", 10*60, true, &DebrickedConfig{})

	files, err := batch.initUpload()

	assert.Len(t, files, 1, "failed to assert that the init deleted one file from the files to be uploaded")
	assert.NoError(t, err)
	assert.Equal(t, 1, batch.ciUploadId)
}

func TestGetDebrickedConfig(t *testing.T) {
	config := GetDebrickedConfig(filepath.Join("testdata", "debricked-config.yaml"))
	configJSON, err := json.Marshal(config)
	assert.Nil(t, err)
	expectedJSON, err := json.Marshal(DebrickedConfig{
		Overrides: []purlConfig{
			{
				PackageURL:  "pkg:npm/lodash",
				Version:     boolOrString{Version: "1.0.0", HasVersion: true},
				FileRegexes: []string{".*/lodash/.*"},
			},
			{
				PackageURL:  "pkg:maven/org.openjfx/javafx-base",
				Version:     boolOrString{Version: "", HasVersion: false},
				FileRegexes: []string{"subpath/org.openjfx/.*"},
			},
		},
	})
	assert.Nil(t, err)
	assert.JSONEq(t, string(configJSON), string(expectedJSON))
}

func TestGetDebrickedConfigUnmarshalError(t *testing.T) {
	config := GetDebrickedConfig(filepath.Join("testdata", "debricked-config-error.yaml"))
	configJSON, err := json.Marshal(config)
	assert.Nil(t, err)
	expectedJSON, err := json.Marshal(nil)
	assert.Nil(t, err)
	assert.JSONEq(t, string(expectedJSON), string(configJSON))
}

func TestGetDebrickedConfigFailure(t *testing.T) {
	config := GetDebrickedConfig("")
	configJSON, err := json.Marshal(config)
	assert.Nil(t, err)
	expectedJSON, err := json.Marshal(nil)
	assert.Nil(t, err)
	assert.JSONEq(t, string(expectedJSON), string(configJSON))
}

func TestMarshalJSONDebrickedConfig(t *testing.T) {
	config, err := json.Marshal(DebrickedConfig{
		Overrides: []purlConfig{
			{
				PackageURL:  "pkg:npm/lodash",
				Version:     boolOrString{Version: "1.0.0", HasVersion: true},
				FileRegexes: []string{".*/lodash/.*"},
			},
			{
				PackageURL:  "pkg:maven/org.openjfx/javafx-base",
				Version:     boolOrString{Version: "", HasVersion: false},
				FileRegexes: []string{"subpath/org.openjfx/.*"},
			},
		},
	})
	expectedJSON := "{\"overrides\":[{\"pURL\":\"pkg:npm/lodash\",\"version\":\"1.0.0\",\"fileRegexes\":[\".*/lodash/.*\"]},{\"pURL\":\"pkg:maven/org.openjfx/javafx-base\",\"version\":false,\"fileRegexes\":[\"subpath/org.openjfx/.*\"]}]}"
	assert.Nil(t, err)
	assert.Equal(t, []byte(expectedJSON), config)
}
