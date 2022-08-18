package uploader

import (
	"testing"
)

func TestNewUploadResult(t *testing.T) {
	status := &uploadStatus{
		Progress:                       100,
		VulnerabilitiesFound:           0,
		UnaffectedVulnerabilitiesFound: 0,
		AutomationsAction:              "",
		AutomationRules:                nil,
		DetailsUrl:                     "",
	}
	result := newUploadResult(status)
	if result == nil {
		t.Error("failed to assert that result was not nil")
	}
}
