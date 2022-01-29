package automation

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

func capturePrintOutput(rule Rule) string {
	var buf bytes.Buffer
	log.SetOutput(&buf)
	rule.Print(&buf)
	log.SetOutput(os.Stderr)

	return buf.String()
}

func TestPrintPassingRule(t *testing.T) {
	rule := Rule{
		RuleDescription: "Rule description",
		RuleActions:     []string{"warnPipeline"},
		RuleLink:        "link",
		HasCves:         false,
		Triggered:       false,
		TriggerEvents:   nil,
	}
	output := capturePrintOutput(rule)

	if len(output) == 0 {
		t.Error("failed to fetch output")
	}
	if !strings.Contains(output, "✔") {
		t.Error("failed to assert that rule did not trigger")
	}

	if !strings.Contains(output, "Rule description") {
		t.Error("failed to assert rule description")
	}

	if !strings.Contains(output, "Manage rule: link") {
		t.Error("failed to assert rule link")
	}
}

func TestPrintTriggeredRule(t *testing.T) {
	rule := Rule{
		RuleDescription: "Rule description",
		RuleActions:     []string{"warnPipeline"},
		RuleLink:        "link",
		HasCves:         false,
		Triggered:       true,
		TriggerEvents:   nil,
	}
	output := capturePrintOutput(rule)

	if len(output) == 0 {
		t.Error("failed to fetch output")
	}
	if strings.Contains(output, "✔") {
		t.Error("failed to assert that rule did trigger")
	}

	if !strings.Contains(output, "Rule description") {
		t.Error("failed to assert rule description")
	}

	if !strings.Contains(output, "Manage rule: link") {
		t.Error("failed to assert rule link")
	}
}

func TestPrintFailingRule(t *testing.T) {
	rule := Rule{
		RuleDescription: "Rule description",
		RuleActions:     []string{"failPipeline"},
		RuleLink:        "link",
		HasCves:         false,
		Triggered:       true,
		TriggerEvents:   nil,
	}
	output := capturePrintOutput(rule)

	if len(output) == 0 {
		t.Error("failed to fetch output")
	}
	if !strings.Contains(output, "⨯") {
		t.Error("failed to assert that rule trigger failure")
	}

	if !strings.Contains(output, "Rule description") {
		t.Error("failed to assert rule description")
	}

	if !strings.Contains(output, "Manage rule: link") {
		t.Error("failed to assert rule link")
	}
}

func TestFailPipeline(t *testing.T) {
	rule := Rule{
		RuleDescription: "Rule description",
		RuleActions:     []string{"warnPipeline", "email"},
		RuleLink:        "link",
		HasCves:         false,
		Triggered:       true,
		TriggerEvents:   nil,
	}
	if rule.FailPipeline() {
		t.Error("failed to assert that rule passed")
	}

	rule.RuleActions = append(rule.RuleActions, "failPipeline")
	if !rule.FailPipeline() {
		t.Error("failed to assert that rule failed")
	}
}
