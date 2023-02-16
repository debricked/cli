package tui

import (
	"bytes"
	"github.com/debricked/cli/pkg/automation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func capturePrintOutput(rule automation.Rule) string {
	var buf bytes.Buffer
	ruleCard := NewRuleCard(&buf, rule)
	ruleCard.Render()

	return buf.String()
}

func TestRenderPassingRule(t *testing.T) {
	rule := automation.Rule{
		RuleDescription: "Rule description",
		RuleActions:     []string{"warnPipeline"},
		RuleLink:        "link",
		HasCves:         false,
		Triggered:       false,
		TriggerEvents:   nil,
	}
	output := capturePrintOutput(rule)

	assert.NotEmpty(t, output, "failed to fetch output")
	assert.Contains(t, output, "✔")
	assert.Contains(t, output, "Rule description")
	assert.Contains(t, output, "Manage rule: link")
}

func TestRenderTriggeredRule(t *testing.T) {
	rule := automation.Rule{
		RuleDescription: "Rule description",
		RuleActions:     []string{"warnPipeline"},
		RuleLink:        "link",
		HasCves:         false,
		Triggered:       true,
		TriggerEvents:   nil,
	}
	output := capturePrintOutput(rule)

	assert.NotEmpty(t, output, "failed to fetch output")
	assert.NotContains(t, output, "✔")
	assert.Contains(t, output, "Rule description")
	assert.Contains(t, output, "Manage rule: link")
}

func TestRenderFailingRule(t *testing.T) {
	rule := automation.Rule{
		RuleDescription: "Rule description",
		RuleActions:     []string{"failPipeline"},
		RuleLink:        "link",
		HasCves:         false,
		Triggered:       true,
		TriggerEvents:   nil,
	}
	output := capturePrintOutput(rule)

	assert.NotEmpty(t, output, "failed to fetch output")
	assert.Contains(t, output, "⨯")
	assert.Contains(t, output, "Rule description")
	assert.Contains(t, output, "Manage rule: link")
}
