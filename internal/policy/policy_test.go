package policy

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	clientTestdata "github.com/debricked/cli/internal/client/testdata"
	internalIO "github.com/debricked/cli/internal/io"
	"github.com/debricked/cli/internal/policy/testdata"
	"github.com/stretchr/testify/assert"
)

const validPolicyPath = "testdata/debricked_policy.json"

func mockResponse(statusCode int, body string) clientTestdata.MockResponse {
	return clientTestdata.MockResponse{
		StatusCode:   statusCode,
		ResponseBody: io.NopCloser(strings.NewReader(body)),
	}
}

func newValidator(response clientTestdata.MockResponse) Validator {
	debClientMock := clientTestdata.NewDebClientMock()
	debClientMock.AddMockResponse(response)

	return Validator{DebClient: debClientMock, FileSystem: internalIO.FileSystem{}}
}

func TestValidateValidPolicyFile(t *testing.T) {
	validator := newValidator(mockResponse(http.StatusOK, "[]"))

	result, err := validator.Validate(validPolicyPath)

	assert.NoError(t, err)
	assert.True(t, result.Valid)
	assert.Equal(t, validPolicyPath, result.File)
	assert.Empty(t, result.Errors)
	assert.NotNil(t, result.Errors, "failed to assert that errors were serializable as an empty array")
}

func TestValidateInvalidPolicyFile(t *testing.T) {
	body := `[
		{"message": "The property name is required", "context": {"property": "policies[0].name"}},
		{"message": "Array must have at least 1 item"}
	]`
	validator := newValidator(mockResponse(http.StatusOK, body))

	result, err := validator.Validate(validPolicyPath)

	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Len(t, result.Errors, 2)
	assert.Equal(t, "policies[0].name", result.Errors[0].PropertyPath)
	assert.Equal(t, "policies[0].name: The property name is required", result.Errors[0].String())
	assert.Equal(t, "Array must have at least 1 item", result.Errors[1].String())
}

func TestValidateRequest(t *testing.T) {
	debClientMock := testdata.NewRecordingDebClientMock()
	debClientMock.AddMockResponse(mockResponse(http.StatusOK, "[]"))
	validator := Validator{DebClient: debClientMock, FileSystem: internalIO.FileSystem{}}

	_, err := validator.Validate(validPolicyPath)

	assert.NoError(t, err)
	assert.Equal(t, "/api/1.0/open/debricked-policy/validate", debClientMock.Uri)

	var request map[string]string
	assert.NoError(t, json.Unmarshal([]byte(debClientMock.Body), &request))
	content, err := os.ReadFile(validPolicyPath)
	assert.NoError(t, err)
	assert.Equal(t, string(content), request["content"])
}

func TestValidateDirectoryPath(t *testing.T) {
	validator := newValidator(mockResponse(http.StatusOK, "[]"))

	result, err := validator.Validate("testdata")

	assert.NoError(t, err)
	assert.True(t, result.Valid)
	assert.Equal(t, filepath.Join("testdata", DefaultFileName), result.File)
}

func TestValidateDefaultPath(t *testing.T) {
	workingDirectory, err := os.Getwd()
	assert.NoError(t, err)
	assert.NoError(t, os.Chdir("testdata"))
	defer func() {
		assert.NoError(t, os.Chdir(workingDirectory))
	}()
	validator := newValidator(mockResponse(http.StatusOK, "[]"))

	result, err := validator.Validate("")

	assert.NoError(t, err)
	assert.True(t, result.Valid)
	assert.Equal(t, DefaultFileName, result.File)
}

func TestValidateMissingPolicyFile(t *testing.T) {
	validator := newValidator(mockResponse(http.StatusOK, "[]"))

	_, err := validator.Validate(filepath.Join("testdata", "missing_policy.json"))

	assert.ErrorIs(t, err, ErrNoPolicyFile)
	assert.ErrorContains(t, err, "missing_policy.json")
}

func TestValidateMissingDefaultPolicyFileInDirectory(t *testing.T) {
	validator := newValidator(mockResponse(http.StatusOK, "[]"))

	_, err := validator.Validate("..")

	assert.ErrorIs(t, err, ErrNoPolicyFile)
	assert.ErrorContains(t, err, DefaultFileName)
}

func TestValidateStatError(t *testing.T) {
	statErr := errors.New("stat-error")
	validator := newValidator(mockResponse(http.StatusOK, "[]"))
	validator.FileSystem = testdata.FileSystemMock{StatError: statErr}

	_, err := validator.Validate(validPolicyPath)

	assert.ErrorContains(t, err, "failed to read")
	assert.ErrorContains(t, err, statErr.Error())
	assert.NotErrorIs(t, err, ErrNoPolicyFile)
}

func TestValidateReadFileError(t *testing.T) {
	readErr := errors.New("read-error")
	validator := newValidator(mockResponse(http.StatusOK, "[]"))
	validator.FileSystem = testdata.FileSystemMock{ReadFileError: readErr}

	_, err := validator.Validate(validPolicyPath)

	assert.ErrorContains(t, err, "failed to read")
	assert.ErrorContains(t, err, readErr.Error())
}

func TestValidateMalformedJson(t *testing.T) {
	validator := newValidator(mockResponse(http.StatusOK, "[]"))

	result, err := validator.Validate(filepath.Join("testdata", "malformed_policy.json"))

	assert.NoError(t, err)
	assert.False(t, result.Valid)
	assert.Len(t, result.Errors, 1)
	assert.Contains(t, result.Errors[0].Message, "invalid JSON on line 5")
}

func TestValidateClientError(t *testing.T) {
	clientErr := errors.New("client-error")
	debClientMock := clientTestdata.NewDebClientMock()
	debClientMock.AddMockResponse(clientTestdata.MockResponse{Error: clientErr})
	validator := Validator{DebClient: debClientMock, FileSystem: internalIO.FileSystem{}}

	_, err := validator.Validate(validPolicyPath)

	assert.ErrorIs(t, err, clientErr)
}

func TestValidateBadStatusCode(t *testing.T) {
	validator := newValidator(mockResponse(http.StatusBadRequest, `{"error": "content is required"}`))

	_, err := validator.Validate(validPolicyPath)

	assert.ErrorContains(t, err, "failed to validate policy file")
	assert.ErrorContains(t, err, "400")
	assert.ErrorContains(t, err, "content is required")
}

func TestValidateMalformedResponse(t *testing.T) {
	validator := newValidator(mockResponse(http.StatusOK, "not json"))

	_, err := validator.Validate(validPolicyPath)

	assert.ErrorContains(t, err, "failed to parse validation response")
}
