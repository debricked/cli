package scan

import (
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
	status, err := newScanStatus(res)
	if err == nil {
		t.Error("failed to assert that error occurred")
	}
	if status != nil {
		t.Error("failed to assert that scanStatus was nil")
	}
}
