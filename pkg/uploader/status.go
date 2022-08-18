package uploader

import (
	"debricked/pkg/automation"
	"encoding/json"
	"io"
	"net/http"
)

type uploadStatus struct {
	Progress                       int               `json:"progress"`
	VulnerabilitiesFound           int               `json:"vulnerabilitiesFound"`
	UnaffectedVulnerabilitiesFound int               `json:"unaffectedVulnerabilitiesFound"`
	AutomationsAction              string            `json:"automationsAction"`
	AutomationRules                []automation.Rule `json:"automationRules"`
	DetailsUrl                     string            `json:"detailsUrl"`
}

func newUploadStatus(response *http.Response) (*uploadStatus, error) {
	status := uploadStatus{}
	data, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	err := json.Unmarshal(data, &status)
	if err != nil {
		return nil, err
	}

	return &status, err
}
