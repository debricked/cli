package tui

import (
	"bytes"
	"fmt"
	"github.com/debricked/cli/pkg/automation"
	"github.com/fatih/color"
	"github.com/jedib0t/go-pretty/v6/list"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"io"
	"strings"
)

const rowLength = 90

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

	style := table.StyleRounded
	style.Format.Footer = text.FormatDefault
	t.SetStyle(style)

	rc.addHeader(t.(*table.Table))
	rc.addDescription(t.(*table.Table))
	rc.addTriggers(t.(*table.Table))

	t.Render()

	fmt.Println()
}

func (rc RuleCard) addDescription(t *table.Table) {
	var descriptionBytes []byte

	description := strings.TrimSuffix(rc.rule.RuleDescription, "\n")
	words := strings.Fields(description)
	length := 0
	for _, word := range words {
		wordWithSpace := word + " "
		wordLen := len(wordWithSpace)
		if length+wordLen < rowLength {
			length = length + wordLen
			descriptionBytes = append(descriptionBytes, []byte(wordWithSpace)...)
		} else {
			length = wordLen
			descriptionBytes = append(descriptionBytes, []byte("\n"+wordWithSpace)...)
		}
	}

	descriptionBytes = append(descriptionBytes, []byte("\n\n")...)
	descriptionBytes = append(descriptionBytes, []byte(fmt.Sprintf("Manage rule: %s", rc.rule.RuleLink))...)

	row := string(descriptionBytes)
	t.AppendRow(table.Row{"", row})
}

func (rc RuleCard) addTriggers(t *table.Table) {
	dependencies := makeDependenciesFromTriggers(rc.rule.TriggerEvents)
	for _, dep := range dependencies {

		var listBuffer bytes.Buffer

		title := fmt.Sprintf("%s:\n", color.BlueString(dep.name))
		underlining := fmt.Sprintf(strings.Repeat("-", len(title)+1) + "\n")
		listBuffer.Write([]byte(title))
		listBuffer.Write([]byte(underlining))

		listWriter := list.NewWriter()

		listWriter.AppendItem(fmt.Sprintf("URL: %s", dep.url))

		rc.addVulnerabilities(listWriter.(*list.List), dep.vulnerabilities)
		rc.addLicenses(listWriter.(*list.List), dep.licenses)

		listWriter.SetOutputMirror(&listBuffer)
		listWriter.Render()

		row := table.Row{"", listBuffer.String()}
		t.AppendFooter(row)
	}
}

func (rc RuleCard) addVulnerabilities(l *list.List, vulnerabilities map[string]vulnerability) {
	if len(vulnerabilities) > 0 {
		l.AppendItem("Vulnerabilities:")
		l.Indent()
		for _, vuln := range vulnerabilities {
			l.UnIndent()
			l.AppendItem(color.YellowString(vuln.name))
			l.Indent()
			l.AppendItem(fmt.Sprintf("URL: %s", vuln.url))
			l.AppendItem(fmt.Sprintf("CVSS2: %s", color.YellowString(vuln.cvss2)))
			l.AppendItem(fmt.Sprintf("CVSS3: %s", color.YellowString(vuln.cvss3)))
			l.UnIndent()
		}
	}
}

func (rc RuleCard) addLicenses(l *list.List, licenses map[string]bool) {
	if len(licenses) > 0 {
		l.AppendItem("Licenses:")
		l.Indent()
		for license := range licenses {
			l.AppendItem(color.YellowString(license))
		}
		l.UnIndent()
	}
}

func (rc RuleCard) addHeader(t *table.Table) {
	triggerStatus := color.GreenString("✔")
	if rc.rule.Triggered {
		triggerStatus = color.YellowString("!")
		if rc.rule.FailPipeline() {
			triggerStatus = color.RedString("⨯")
		}
	}
	t.AppendHeader(table.Row{triggerStatus, "Automation Rule"})
}
