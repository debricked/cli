package automation

const failPipeline = "failPipeline"

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
	for _, action := range rule.RuleActions {
		if action == failPipeline {
			return true
		}
	}

	return false
}
