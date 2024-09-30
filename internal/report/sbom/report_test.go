package sbom

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/debricked/cli/internal/client/testdata"
	ioTestData "github.com/debricked/cli/internal/io/testdata"
	"github.com/stretchr/testify/assert"
)

func TestOrderError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusOK})
	reporter := Reporter{DebClient: debClientMock, FileWriter: &ioTestData.FileWriterMock{}}
	err := reporter.Order(OrderArgs{CommitID: "", RepositoryID: ""})
	assert.Error(t, err)
}

func TestOrder(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: io.NopCloser(strings.NewReader("{}")),
	})
	debClientMock.AddMockResponse(testdata.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: io.NopCloser(strings.NewReader("{}")),
	})
	reporter := Reporter{DebClient: debClientMock, FileWriter: &ioTestData.FileWriterMock{}}
	err := reporter.Order(OrderArgs{CommitID: "", RepositoryID: ""})
	assert.NoError(t, err)
}

func TestOrderDownloadErr(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: io.NopCloser(strings.NewReader("{}")),
	})
	debClientMock.AddMockResponse(testdata.MockResponse{
		StatusCode: http.StatusForbidden,
	})
	reporter := Reporter{DebClient: debClientMock, FileWriter: &ioTestData.FileWriterMock{}}
	err := reporter.Order(OrderArgs{CommitID: "", RepositoryID: ""})
	assert.Error(t, err)
}

func TestOrderArgsError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusOK})
	reporter := Reporter{DebClient: debClientMock}
	err := reporter.Order("")
	assert.Error(t, err)
}

func TestGenerateOK(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: io.NopCloser(strings.NewReader("{}")),
	})
	reporter := Reporter{DebClient: debClientMock, FileWriter: &ioTestData.FileWriterMock{}}
	uuid, err := reporter.generate(orderArgs())
	assert.NoError(t, err)
	assert.NotNil(t, uuid)
}

func TestGenerateSubscriptionError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{
		StatusCode:   http.StatusPaymentRequired,
		ResponseBody: io.NopCloser(strings.NewReader("{}")),
	})
	reporter := Reporter{DebClient: debClientMock, FileWriter: &ioTestData.FileWriterMock{}}
	uuid, err := reporter.generate(orderArgs())
	assert.Error(t, err)
	assert.NotNil(t, uuid)
}

func TestGenerateError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{
		StatusCode:   http.StatusForbidden,
		ResponseBody: io.NopCloser(strings.NewReader("{}")),
	})
	reporter := Reporter{DebClient: debClientMock, FileWriter: &ioTestData.FileWriterMock{}}
	uuid, err := reporter.generate(orderArgs())
	assert.Error(t, err)
	assert.NotNil(t, uuid)
}

func TestGenerateDefaultGetError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{Error: errors.New("")})
	reporter := Reporter{DebClient: debClientMock, FileWriter: &ioTestData.FileWriterMock{}}
	res, err := reporter.generate(orderArgs())
	assert.Error(t, err)
	assert.Equal(t, "", res)
}

func TestDownloadOK(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusOK})
	reporter := Reporter{DebClient: debClientMock, FileWriter: &ioTestData.FileWriterMock{}}
	res, err := reporter.download("")
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestDownloadTooLongQueue(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusCreated})
	reporter := Reporter{DebClient: debClientMock, FileWriter: &ioTestData.FileWriterMock{}}
	res, err := reporter.download("")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestDownloadDefaultError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusForbidden})
	reporter := Reporter{DebClient: debClientMock, FileWriter: &ioTestData.FileWriterMock{}}
	res, err := reporter.download("")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestDownloadDefaultGetError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{Error: errors.New("")})
	reporter := Reporter{DebClient: debClientMock, FileWriter: &ioTestData.FileWriterMock{}}
	res, err := reporter.download("")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestParseURL(t *testing.T) {
	testURL := "https://debricked.com/app/en/repository/0/commit/1"
	clientMock := testdata.NewDebClientMock()
	reporter := Reporter{DebClient: clientMock, FileWriter: &ioTestData.FileWriterMock{}}
	repositoryID, commitID, err := reporter.ParseDetailsURL(testURL)

	assert.NoError(t, err)
	assert.Equal(t, repositoryID, "0")
	assert.Equal(t, commitID, "1")
}

func TestParseURLFormatErr(t *testing.T) {
	testURL := "https://debricked.com/app/en/repository/0"
	clientMock := testdata.NewDebClientMock()
	reporter := Reporter{DebClient: clientMock, FileWriter: &ioTestData.FileWriterMock{}}
	_, _, err := reporter.ParseDetailsURL(testURL)

	assert.Error(t, err)
}

func TestWriteSBOM(t *testing.T) {
	clientMock := testdata.NewDebClientMock()
	fileWriter := &ioTestData.FileWriterMock{
		CreateErr: errors.New(""),
	}
	reporter := Reporter{DebClient: clientMock, FileWriter: fileWriter}
	err := reporter.writeSBOM("", "", "", nil)
	assert.Error(t, err)
}

func orderArgs() OrderArgs {
	return OrderArgs{
		Vulnerabilities: false,
		Licenses:        false,
		Branch:          "",
		CommitID:        "",
		RepositoryID:    "",
		Output:          "",
	}
}
