package sbom

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/debricked/cli/internal/client/testdata"
	"github.com/stretchr/testify/assert"
)

func TestOrderError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusOK})
	reporter := Reporter{DebClient: debClientMock}
	args := OrderArgs{CommitID: "", RepositoryID: ""}
	err := reporter.Order(args)
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
	reporter := Reporter{DebClient: debClientMock}
	args := OrderArgs{CommitID: "", RepositoryID: ""}
	err := reporter.Order(args)
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
	reporter := Reporter{DebClient: debClientMock}
	args := OrderArgs{CommitID: "", RepositoryID: ""}
	err := reporter.Order(args)
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
	reporter := Reporter{DebClient: debClientMock}
	uuid, err := reporter.generate("", "", "", false, false)
	assert.NoError(t, err)
	assert.NotNil(t, uuid)
}

func TestGenerateSubscriptionError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{
		StatusCode:   http.StatusPaymentRequired,
		ResponseBody: io.NopCloser(strings.NewReader("{}")),
	})
	reporter := Reporter{DebClient: debClientMock}
	uuid, err := reporter.generate("", "", "", false, false)
	assert.Error(t, err)
	assert.NotNil(t, uuid)
}

func TestGenerateError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{
		StatusCode:   http.StatusForbidden,
		ResponseBody: io.NopCloser(strings.NewReader("{}")),
	})
	reporter := Reporter{DebClient: debClientMock}
	uuid, err := reporter.generate("", "", "", false, false)
	assert.Error(t, err)
	assert.NotNil(t, uuid)
}

func TestDownloadOK(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusOK})
	reporter := Reporter{DebClient: debClientMock}
	res, err := reporter.download("")
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestDownloadTooLongQueue(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusCreated})
	reporter := Reporter{DebClient: debClientMock}
	res, err := reporter.download("")
	assert.Error(t, err)
	assert.NotNil(t, res)
}

func TestDownloadDefaultError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusForbidden})
	reporter := Reporter{DebClient: debClientMock}
	res, err := reporter.download("")
	assert.Error(t, err)
	assert.NotNil(t, res)
}

func TestParseURL(t *testing.T) {
	testURL := "https://debricked.com/app/en/repository/0/commit/1"
	clientMock := testdata.NewDebClientMock()
	reporter := Reporter{DebClient: clientMock}
	repositoryID, commitID, err := reporter.ParseDetailsURL(testURL)

	assert.NoError(t, err)
	assert.Equal(t, repositoryID, "0")
	assert.Equal(t, commitID, "1")
}

func TestParseURLFormatErr(t *testing.T) {
	testURL := "https://debricked.com/app/en/repository/0"
	clientMock := testdata.NewDebClientMock()
	reporter := Reporter{DebClient: clientMock}
	_, _, err := reporter.ParseDetailsURL(testURL)

	assert.Error(t, err)
}
