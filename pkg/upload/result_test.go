package upload

import (
	"github.com/stretchr/testify/assert"
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

	assert.NotNil(t, result)
}
