package upload

import (
	"github.com/debricked/cli/internal/automation"
)

type UploadResult struct {
	VulnerabilitiesFound           int               `json:"vulnerabilitiesFound"`
	UnaffectedVulnerabilitiesFound int               `json:"unaffectedVulnerabilitiesFound"`
	AutomationsAction              string            `json:"automationsAction"`
	AutomationRules                []automation.Rule `json:"automationRules"`
	DetailsUrl                     string            `json:"detailsUrl"`
}

func newUploadResult(status *uploadStatus) *UploadResult {
	return &UploadResult{
		status.VulnerabilitiesFound,
		status.UnaffectedVulnerabilitiesFound,
		status.AutomationsAction,
		status.AutomationRules,
		status.DetailsUrl,
	}
}
