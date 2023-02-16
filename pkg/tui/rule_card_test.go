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

func TestRenderTriggerEvents(t *testing.T) {
	rule := automation.Rule{
		RuleDescription: "Rule description Rule description Rule description Rule description Rule description Rule description Rule description Rule description ",
		RuleActions:     []string{"failPipeline"},
		RuleLink:        "link",
		HasCves:         false,
		Triggered:       true,
	}

	rule.TriggerEvents = []automation.TriggerEvent{
		{
			Dependency:     "dependency-1",
			DependencyLink: "dependency-1-url",
			Licenses:       []string{"MIT"},
			Cve:            "CVE-1",
			Cvss2:          5.0,
			Cvss3:          5.5,
			CveLink:        "CVE-1-URL",
		},
		{
			Dependency:     "dependency-1",
			DependencyLink: "dependency-1-url",
			Licenses:       []string{"MIT"},
			Cve:            "CVE-2",
			Cvss2:          6.0,
			Cvss3:          6.5,
			CveLink:        "CVE-2-URL",
		},
	}

	output := capturePrintOutput(rule)

	assert.NotEmpty(t, output, "failed to fetch output")
	assert.Contains(t, output, "⨯")
	assert.Contains(t, output, "Rule description")
	assert.Contains(t, output, "Manage rule: link")

	assert.Contains(t, output, "\u001B[34mdependency-1\u001B[0m:")
	assert.Contains(t, output, "* URL: dependency-1-url")
	assert.Contains(t, output, "* \u001B[33mCVE-1\u001B[0m")
	assert.Contains(t, output, "* URL: CVE-1-URL")
	assert.Contains(t, output, "* CVSS2: \u001B[33m5\u001B[0m")
	assert.Contains(t, output, "* CVSS3: \u001B[33m5.5\u001B[0m")
	assert.Contains(t, output, "* \u001B[33mCVE-2\u001B[0m")
	assert.Contains(t, output, "URL: CVE-2-URL")
	assert.Contains(t, output, "* \u001B[33mCVE-2\u001B[0m")
	assert.Contains(t, output, "* CVSS2: \u001B[33m6\u001B[0m")
	assert.Contains(t, output, "* CVSS3: \u001B[33m6.5\u001B[0m")
	assert.Contains(t, output, "* Licenses:")
	assert.Contains(t, output, "\u001B[33mMIT\u001B[0m")
}
