package automation

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"io"
)

type Rule struct {
	RuleDescription string         `json:"ruleDescription"`
	RuleActions     []string       `json:"ruleActions"`
	RuleLink        string         `json:"ruleLink"`
	HasCves         bool           `json:"hasCves"`
	Triggered       bool           `json:"triggered"`
	TriggerEvents   []TriggerEvent `json:"triggerEvents"`
}

func (rule *Rule) Print(mirror io.Writer) {
	t := table.NewWriter()
	t.SetOutputMirror(mirror)
	triggerStatus := color.GreenString("✔")
	if rule.Triggered {
		triggerStatus = color.YellowString("!")
		failPipeline := false
		for _, action := range rule.RuleActions {
			if action == "failPipeline" {
				failPipeline = true

				break
			}
		}
		if failPipeline {
			triggerStatus = color.RedString("⨯")
		}
	}
	t.AppendHeader(table.Row{triggerStatus, "Automation Rule"})
	t.AppendRow(table.Row{"", rule.RuleDescription + "\n\n" + fmt.Sprintf("Manage rule: %s", rule.RuleLink)})
	t.SetStyle(table.StyleLight)
	t.Render()
	fmt.Println()
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
