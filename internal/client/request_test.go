package client

import (
	"bytes"
	"net/http"
	"testing"

	testdataClient "github.com/debricked/cli/internal/client/testdata/client"
	"github.com/stretchr/testify/assert"
)

func TestGetNilRes(t *testing.T) {
	clientMock := testdataClient.NewMock()
	debClient := NewDebClient(nil, clientMock)

	response, err := get("", debClient, true, "") //nolint:bodyclose

	assert.ErrorIs(t, NoResErr, err)
	assert.Nil(t, response)
}

func TestPostNilRes(t *testing.T) {
	clientMock := testdataClient.NewMock()
	debClient := NewDebClient(nil, clientMock)

	response, err := post("", debClient, "application/json", bytes.NewBuffer(nil), true) //nolint:bodyclose
	assert.ErrorIs(t, NoResErr, err)
	assert.Nil(t, response)
}

func TestInterpretNilRes(t *testing.T) {
	response, err := interpret(nil, func() (*http.Response, error) { //nolint:bodyclose
		return nil, NoResErr
	}, nil, true)

	assert.Nil(t, response)
	assert.ErrorIs(t, NoResErr, err)
}
