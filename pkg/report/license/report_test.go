package license

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/debricked/cli/pkg/client/testdata"
	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	addCommitIdMockResponse(debClientMock)
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusOK})
	reporter := Reporter{DebClient: debClientMock}
	args := OrderArgs{Email: ""}

	err := reporter.Order(args)

	assert.NoError(t, err)
}

func TestOrderBadArgs(t *testing.T) {
	debClientMock := &testdata.DebClientMock{}
	reporter := Reporter{DebClient: debClientMock}
	args := struct{}{}

	err := reporter.Order(args)

	assert.ErrorIs(t, err, ArgsError)
}

func TestGetCommitError(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	errorAssertion := errors.New("commitError")
	debClientMock.AddMockResponse(testdata.MockResponse{Error: errorAssertion})
	reporter := Reporter{DebClient: debClientMock}
	args := OrderArgs{Email: ""}

	err := reporter.Order(args)

	assert.ErrorIs(t, err, errorAssertion)

}

func TestOrderUnauthorized(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	errorAssertion := errors.New("unauthorized")
	debClientMock.AddMockResponse(testdata.MockResponse{Error: errorAssertion})
	reporter := Reporter{DebClient: debClientMock}
	args := OrderArgs{Email: ""}

	err := reporter.Order(args)

	assert.ErrorIs(t, err, errorAssertion)
}

func TestOrderForbidden(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	addCommitIdMockResponse(debClientMock)
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusForbidden})
	reporter := Reporter{DebClient: debClientMock}
	args := OrderArgs{Email: ""}

	err := reporter.Order(args)

	assert.ErrorIs(t, err, SubscriptionError)
}

func TestOrderNotOkResponse(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	addCommitIdMockResponse(debClientMock)
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusTeapot})
	reporter := Reporter{DebClient: debClientMock}
	args := OrderArgs{Email: ""}

	err := reporter.Order(args)

	assert.ErrorContains(t, err, fmt.Sprintf("failed to order report. Status code: %d", http.StatusTeapot))
}

func TestGetCommitId(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	addCommitIdMockResponse(debClientMock)
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusOK})
	reporter := Reporter{DebClient: debClientMock}
	args := OrderArgs{Email: ""}

	err := reporter.Order(args)

	assert.NoError(t, err)
}

func TestGetCommitIdUnauthorized(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	errorAssertion := errors.New("unauthorized")
	debClientMock.AddMockResponse(testdata.MockResponse{Error: errorAssertion})
	reporter := Reporter{DebClient: debClientMock}

	_, err := reporter.getCommitId("")

	assert.ErrorIs(t, err, errorAssertion)
}

func TestGetCommitIdForbidden(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusForbidden})
	reporter := Reporter{DebClient: debClientMock}

	_, err := reporter.getCommitId("")

	assert.ErrorIs(t, err, SubscriptionError)
}

func TestGetCommitIdNotOkResponse(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{StatusCode: http.StatusTeapot})
	reporter := Reporter{DebClient: debClientMock}

	_, err := reporter.getCommitId("")

	assert.ErrorContains(t, err, "no commit was found with the name")
}

func TestGetCommitIdNoResult(t *testing.T) {
	debClientMock := testdata.NewDebClientMock()
	debClientMock.AddMockResponse(testdata.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: createIoReadCloserFromCommit(nil),
	})
	reporter := Reporter{DebClient: debClientMock}

	_, err := reporter.getCommitId("")

	assert.ErrorContains(t, err, "no commit was found with the name")
}

func addCommitIdMockResponse(mockClient *testdata.DebClientMock) {
	c := commit{
		FileIds:     []int{},
		Id:          0,
		Name:        "commit-hash",
		ReleaseData: "",
	}
	mockResponse := testdata.MockResponse{
		StatusCode:   http.StatusOK,
		ResponseBody: createIoReadCloserFromCommit(&c),
		Error:        nil,
	}
	mockClient.AddMockResponse(mockResponse)
}

func createIoReadCloserFromCommit(c *commit) io.ReadCloser {
	var commitResponse []commit
	if c != nil {
		commitResponse = append(commitResponse, *c)
	}
	body, _ := json.Marshal(commitResponse)
	reader := bytes.NewReader(body)

	return io.NopCloser(reader)
}
