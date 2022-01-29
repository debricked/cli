package scan

import (
	"debricked/pkg/automation"
	"encoding/json"
	"io"
	"net/http"
)

type scanStatus struct {
	Progress                       int               `json:"progress"`
	VulnerabilitiesFound           int               `json:"vulnerabilitiesFound"`
	UnaffectedVulnerabilitiesFound int               `json:"unaffectedVulnerabilitiesFound"`
	AutomationsAction              string            `json:"automationsAction"`
	AutomationRules                []automation.Rule `json:"automationRules"`
	DetailsUrl                     string            `json:"detailsUrl"`
}

func newScanStatus(response *http.Response) (*scanStatus, error) {
	status := scanStatus{}
	data, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	err := json.Unmarshal(data, &status)
	if err != nil {
		return nil, err
	}

	return &status, err
}
