package upload

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestNewScanStatusBadResponse(t *testing.T) {
	res := &http.Response{
		Status:           "",
		StatusCode:       0,
		Proto:            "",
		ProtoMajor:       0,
		ProtoMinor:       0,
		Header:           nil,
		Body:             http.NoBody,
		ContentLength:    0,
		TransferEncoding: nil,
		Close:            false,
		Uncompressed:     false,
		Trailer:          nil,
		Request:          nil,
		TLS:              nil,
	}
	status, err := newUploadStatus(res)

	assert.Error(t, err)
	assert.Nil(t, status)
}
