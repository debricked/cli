package tui

import (
	"fmt"
	"github.com/debricked/cli/pkg/automation"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"io"
	"strings"
)

const rowLength = 72

type RuleCard struct {
	mirror io.Writer
	rule   automation.Rule
}

func NewRuleCard(mirror io.Writer, rule automation.Rule) RuleCard {
	return RuleCard{mirror: mirror, rule: rule}
}

func (rc RuleCard) Render() {
	t := table.NewWriter()
	t.SetOutputMirror(rc.mirror)
	triggerStatus := color.GreenString("✔")
	if rc.rule.Triggered {
		triggerStatus = color.YellowString("!")
		failPipeline := false
		for _, action := range rc.rule.RuleActions {
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
	rc.addDescription(t.(*table.Table))
	//t.AppendFooter(table.Row{"", "Failed because of X and Y Failed because of X and Y Failed because of X and Y Failed because of X and Y"})
	style := table.StyleRounded
	style.Format.Footer = text.FormatDefault
	t.SetStyle(style)
	t.Render()
	fmt.Println()
}

func (rc RuleCard) addDescription(t *table.Table) {
	var bytes []byte

	description := strings.TrimSuffix(rc.rule.RuleDescription, "\n")
	words := strings.Fields(description)
	length := 0
	for _, word := range words {
		wordWithSpace := word + " "
		wordLen := len(wordWithSpace)
		if length+wordLen < rowLength {
			length = length + wordLen
			bytes = append(bytes, []byte(wordWithSpace)...)
		} else {
			length = wordLen
			bytes = append(bytes, []byte("\n"+wordWithSpace)...)
		}
	}
	bytes = append(bytes, []byte("\n\n")...)
	bytes = append(bytes, []byte(fmt.Sprintf("Manage rule: %s", rc.rule.RuleLink))...)

	row := string(bytes)
	t.AppendRow(table.Row{"", row})
}
func addTriggers(t *table.Table, triggerEvent []automation.TriggerEvent) {

}
