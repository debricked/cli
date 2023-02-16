package automation

type Rule struct {
	RuleDescription string         `json:"ruleDescription"`
	RuleActions     []string       `json:"ruleActions"`
	RuleLink        string         `json:"ruleLink"`
	HasCves         bool           `json:"hasCves"`
	Triggered       bool           `json:"triggered"`
	TriggerEvents   []TriggerEvent `json:"triggerEvents"`
}

// FailPipeline checks if rule should fail the pipeline
func (rule *Rule) FailPipeline() bool {
	if rule.Triggered {
		for _, action := range rule.RuleActions {
			if action == "failPipeline" {
				return true
			}
		}
	}

	return false
}

type TriggerEvent struct {
	Dependency     string   `json:"dependency"`
	DependencyLink string   `json:"dependencyLink"`
	Licenses       []string `json:"licenses"`
	Cve            string   `json:"cve"`
	Cvss2          float32  `json:"cvss2"`
	Cvss3          float32  `json:"cvss3"`
	CveLink        string   `json:"cveLink"`
}
