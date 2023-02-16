package automation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFailPipeline(t *testing.T) {
	rule := Rule{
		RuleDescription: "Rule description",
		RuleActions:     []string{"warnPipeline", "email"},
		RuleLink:        "link",
		HasCves:         false,
		Triggered:       true,
		TriggerEvents:   nil,
	}
	assert.False(t, rule.FailPipeline(), "failed to assert that rule passed")

	rule.RuleActions = append(rule.RuleActions, "failPipeline")

	assert.True(t, rule.FailPipeline(), "failed to assert that rule failed")
}

func TestFailPipelineUnTriggered(t *testing.T) {
	rule := Rule{
		RuleDescription: "Rule description",
		RuleActions:     []string{"warnPipeline", "email", "failPipeline"},
		RuleLink:        "link",
		HasCves:         false,
		Triggered:       false,
		TriggerEvents:   nil,
	}
	assert.True(t, rule.FailPipeline(), "failed to assert that rule failed")
}
