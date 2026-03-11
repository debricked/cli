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
	batch := newUploadBatch(&c, groups, metaObj, "CLI", 10*60, true, &DebrickedConfig{}, true, false)
	var buf bytes.Buffer
	log.SetOutput(&buf)
	err = batch.upload()
	log.SetOutput(os.Stderr)
	output := buf.String()

	assert.Empty(t, output)
	assert.ErrorContains(t, err, "Failed to initialize a scan for")
}

func TestInitAnalysisWithoutAnyFiles(t *testing.T) {
	batch := newUploadBatch(nil, file.Groups{}, nil, "CLI", 10*60, true, &DebrickedConfig{}, true, false)
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
	batch := newUploadBatch(&c, groups, metaObj, "CLI", 10*60, true, &DebrickedConfig{}, true, false)

	uploadResult, err := batch.wait()

	assert.True(t, uploadResult.LongQueue)
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
	batch := newUploadBatch(&c, groups, metaObj, "CLI", 10*60, true, &DebrickedConfig{}, true, false)

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
	batch := newUploadBatch(&c, groups, metaObj, "CLI", 10*60, true, &DebrickedConfig{}, true, false)

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
	batch := newUploadBatch(&c, groups, metaObj, "CLI", 10*60, true, &DebrickedConfig{}, true, false)

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

func TestGetDebrickedConfigIgnore(t *testing.T) {
	config := GetDebrickedConfig(filepath.Join("testdata", "debricked-config-ignore.yaml"))
	configJSON, err := json.Marshal(config)
	assert.Nil(t, err)
	expectedJSON, err := json.Marshal(DebrickedConfig{
		Ignore: &IgnoreConfig{
			Packages: []IgnorePackage{
				{PURL: "pkg:npm/verdaccio", Version: "3.7.0"},
				{PURL: "pkg:npm/chart.js"},
				{PURL: "pkg:nuget/simpleinjector", Version: "4.7.1"},
			},
		},
	})
	assert.Nil(t, err)
	assert.JSONEq(t, string(configJSON), string(expectedJSON))
}

func TestGetDebrickedConfigOverridesIgnore(t *testing.T) {
	config := GetDebrickedConfig(filepath.Join("testdata", "debricked-config-override-ignore.yaml"))
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
		Ignore: &IgnoreConfig{
			Packages: []IgnorePackage{
				{PURL: "pkg:npm/verdaccio", Version: "3.7.0"},
				{PURL: "pkg:npm/chart.js"},
				{PURL: "pkg:nuget/simpleinjector", Version: "4.7.1"},
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
	expectedJSON := "{\"override\":[{\"pURL\":\"pkg:npm/lodash\",\"version\":\"1.0.0\",\"fileRegexes\":[\".*/lodash/.*\"]},{\"pURL\":\"pkg:maven/org.openjfx/javafx-base\",\"version\":false,\"fileRegexes\":[\"subpath/org.openjfx/.*\"]}]}"
	assert.Nil(t, err)
	assert.Equal(t, []byte(expectedJSON), config)
}

func TestMarshalJSONDebrickedConfigIgnoreOnly(t *testing.T) {
	config, err := json.Marshal(DebrickedConfig{
		Ignore: &IgnoreConfig{
			Packages: []IgnorePackage{
				{PURL: "pkg:npm/verdaccio", Version: "3.7.0"},
				{PURL: "pkg:npm/chart.js"},
			},
		},
	})
	expectedJSON := "{\"ignore\":{\"packages\":[{\"pURL\":\"pkg:npm/verdaccio\",\"version\":\"3.7.0\"},{\"pURL\":\"pkg:npm/chart.js\"}]}}"
	assert.Nil(t, err)
	assert.Equal(t, []byte(expectedJSON), config)
}

func TestMarshalJSONDebrickedConfigBoth(t *testing.T) {
	config, err := json.Marshal(DebrickedConfig{
		Overrides: []purlConfig{
			{
				PackageURL:  "pkg:npm/lodash",
				Version:     boolOrString{Version: "1.0.0", HasVersion: true},
				FileRegexes: []string{".*/lodash/.*"},
			},
		},
		Ignore: &IgnoreConfig{
			Packages: []IgnorePackage{
				{PURL: "pkg:npm/chart.js"},
			},
		},
	})
	expectedJSON := "{\"override\":[{\"pURL\":\"pkg:npm/lodash\",\"version\":\"1.0.0\",\"fileRegexes\":[\".*/lodash/.*\"]}],\"ignore\":{\"packages\":[{\"pURL\":\"pkg:npm/chart.js\"}]}}"
	assert.Nil(t, err)
	assert.Equal(t, []byte(expectedJSON), config)
}

func TestGetDebrickedConfigSingularOverride(t *testing.T) {
	config := GetDebrickedConfig(filepath.Join("testdata", "debricked-config-singular-override.yaml"))
	configJSON, err := json.Marshal(config)
	assert.Nil(t, err)
	expectedJSON, err := json.Marshal(DebrickedConfig{
		Overrides: []purlConfig{
			{
				PackageURL:  "pkg:npm/lodash",
				Version:     boolOrString{Version: "1.0.0", HasVersion: true},
				FileRegexes: []string{".*/lodash/.*"},
			},
		},
	})
	assert.Nil(t, err)
	assert.JSONEq(t, string(configJSON), string(expectedJSON))
}

func TestGetDebrickedConfigPolicies(t *testing.T) {
	config := GetDebrickedConfig(filepath.Join("testdata", "debricked-config-policies.yaml"))
	configJSON, err := json.Marshal(config)
	assert.Nil(t, err)
	expectedJSON, err := json.Marshal(DebrickedConfig{
		Overrides: []purlConfig{
			{
				PackageURL:  "pkg:npm/lodash",
				Version:     boolOrString{Version: "1.0.0", HasVersion: true},
				FileRegexes: []string{"chart.js-2.6.0.tgz"},
			},
		},
		Ignore: &IgnoreConfig{
			Packages: []IgnorePackage{
				{PURL: "pkg:maven/javax.transaction/jta"},
				{PURL: "pkg:maven/org.quartz-scheduler/quartz"},
				{PURL: "pkg:maven/com.google.guava/guava", Version: "1.1.1"},
				{PURL: "pkg:maven/com.googlecode.json-simple/json-simplea", Version: "1.1.1"},
				{PURL: "pkg:maven/com.fasterxml.jackson.core/jackson-databind"},
			},
		},
		Policies: &PoliciesConfig{
			Allow: &PolicyPackages{
				Packages: []string{
					"pkg:npm/lodash@4.17.21",
					"pkg:maven/org.springframework/spring-core@5.3.20",
					"react",
					"express",
					"axios@1.3.0",
					"lodash@>=4.17.21,<5.0.0",
					"log4j@2.15.0-2.17.1",
				},
			},
			Deny: &PolicyPackages{
				Packages: []string{
					"pkg:npm/request",
					"colors@<=1.4.0",
					"node-ipc@<=9.2.1",
					"pkg:pypi/setuptools@<65.0.0",
					"pkg:pypi/gpl-restricted-package",
					"proprietary-lib",
				},
			},
		},
	})
	assert.Nil(t, err)
	assert.JSONEq(t, string(configJSON), string(expectedJSON))
}

func TestGetDebrickedConfigPoliciesOnly(t *testing.T) {
	config := GetDebrickedConfig(filepath.Join("testdata", "debricked-config-policies-only.yaml"))
	configJSON, err := json.Marshal(config)
	assert.Nil(t, err)
	expectedJSON, err := json.Marshal(DebrickedConfig{
		Policies: &PoliciesConfig{
			Allow: &PolicyPackages{
				Packages: []string{
					// PURL format
					"pkg:npm/lodash@4.17.21",
					"pkg:maven/org.springframework/spring-core@5.3.20",
					"pkg:pypi/requests@2.28.0",
					"pkg:nuget/Newtonsoft.Json@13.0.1",
					"pkg:npm/@angular/core@15.0.0",
					// Name only
					"react",
					"webpack",
					"express",
					"typescript",
					// Name@version
					"axios@1.3.0",
					"lodash@4.17.21",
					"vue@3.2.45",
					// Version ranges
					"django@>=3.2.0,<5.0.0",
					"spring-boot@>=2.7.0,<3.0.0",
					"lodash@>=4.17.21,<5.0.0",
					"log4j@2.15.0-2.17.1",
					// Comparison operators
					"pytest@>=7.0.0",
					"guava@>=31.0.0",
				},
			},
			Deny: &PolicyPackages{
				Packages: []string{
					// PURL format
					"pkg:npm/request",
					"pkg:npm/event-stream@3.3.6",
					"pkg:maven/log4j/log4j@1.2.17",
					"pkg:pypi/pycrypto",
					"pkg:npm/flatmap-stream",
					// Name only
					"colors",
					"node-ipc",
					// Version ranges and constraints
					"setuptools@<65.0.0",
					"minimist@<1.2.6",
					"pillow@<8.3.2",
					"log4j@1.0-2.14.1",
					"commons-collections@<=3.2.1",
					// Comparison operators
					"System.Text.Encodings.Web@<4.7.2",
					"moment@<=2.29.1",
				},
			},
		},
	})
	assert.Nil(t, err)
	assert.JSONEq(t, string(configJSON), string(expectedJSON))
}
