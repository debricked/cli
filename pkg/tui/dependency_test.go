package tui

import (
	"fmt"
	"github.com/debricked/cli/pkg/automation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMakeDependenciesFromTriggersNoTriggers(t *testing.T) {
	deps := makeDependenciesFromTriggers([]automation.TriggerEvent{})

	assert.Empty(t, deps)
}

func TestMakeDependenciesFromTriggersOneDependency(t *testing.T) {
	trigger1 := automation.TriggerEvent{
		Dependency:     "dependency",
		DependencyLink: "",
		Licenses:       nil,
		Cve:            "",
		Cvss2:          0,
		Cvss3:          0,
		CveLink:        "",
	}
	trigger2 := automation.TriggerEvent{
		Dependency:     "dependency",
		DependencyLink: "",
		Licenses:       []string{""},
		Cve:            "",
		Cvss2:          0,
		Cvss3:          0,
		CveLink:        "",
	}
	triggers := []automation.TriggerEvent{trigger1, trigger2}

	deps := makeDependenciesFromTriggers(triggers)

	assert.Len(t, deps, 1)

	dep := deps[trigger1.Dependency]
	assert.Equal(t, trigger1.Dependency, dep.name)
	assert.Equal(t, trigger1.DependencyLink, dep.url)
	assert.Equal(t, trigger2.Dependency, dep.name)
	assert.Equal(t, trigger2.DependencyLink, dep.url)
	assert.Empty(t, dep.vulnerabilities)
	assert.Empty(t, dep.licenses)
}

func TestMakeDependenciesFromTriggersManyDependencies(t *testing.T) {
	trigger1 := automation.TriggerEvent{
		Dependency:     "dependency-1",
		DependencyLink: "",
		Licenses:       []string{"MIT"},
		Cve:            "",
		Cvss2:          0,
		Cvss3:          0,
		CveLink:        "",
	}
	trigger2 := automation.TriggerEvent{
		Dependency:     "dependency-2",
		DependencyLink: "",
		Licenses:       []string{"AGPL"},
		Cve:            "",
		Cvss2:          0,
		Cvss3:          0,
		CveLink:        "",
	}
	triggers := []automation.TriggerEvent{trigger1, trigger2}

	deps := makeDependenciesFromTriggers(triggers)

	assert.Len(t, deps, 2)

	dep1 := deps[trigger1.Dependency]
	assert.Equal(t, trigger1.Dependency, dep1.name)
	assert.Equal(t, trigger1.DependencyLink, dep1.url)
	assert.Empty(t, dep1.vulnerabilities)
	assert.Len(t, dep1.licenses, 1)

	dep2 := deps[trigger2.Dependency]
	assert.Equal(t, trigger2.Dependency, dep2.name)
	assert.Equal(t, trigger2.DependencyLink, dep2.url)
	assert.Empty(t, dep2.vulnerabilities)
	assert.Len(t, dep2.licenses, 1)
}

func TestMakeDependenciesFromTriggersManyDependenciesManyVulns(t *testing.T) {
	trigger1 := automation.TriggerEvent{
		Dependency:     "dependency-1",
		DependencyLink: "",
		Licenses:       []string{""},
		Cve:            "CVE-1",
		Cvss2:          0,
		Cvss3:          0,
		CveLink:        "",
	}
	trigger2 := automation.TriggerEvent{
		Dependency:     "dependency-2",
		DependencyLink: "",
		Licenses:       nil,
		Cve:            "CVE-2",
		Cvss2:          0,
		Cvss3:          0,
		CveLink:        "",
	}
	triggers := []automation.TriggerEvent{trigger1, trigger2}

	deps := makeDependenciesFromTriggers(triggers)

	assert.Len(t, deps, 2)

	dep1 := deps[trigger1.Dependency]
	assert.Equal(t, trigger1.Dependency, dep1.name)
	assert.Equal(t, trigger1.DependencyLink, dep1.url)
	assert.Len(t, dep1.vulnerabilities, 1)
	vuln := dep1.vulnerabilities[trigger1.Cve]
	assert.Equal(t, trigger1.Cve, vuln.name)
	assert.Equal(t, trigger1.CveLink, vuln.url)
	assert.Equal(t, fmt.Sprintf("%g", trigger1.Cvss2), vuln.cvss2)
	assert.Equal(t, fmt.Sprintf("%g", trigger1.Cvss3), vuln.cvss3)
	assert.Empty(t, dep1.licenses)

	dep2 := deps[trigger2.Dependency]
	assert.Equal(t, trigger2.Dependency, dep2.name)
	assert.Equal(t, trigger2.DependencyLink, dep2.url)
	assert.Len(t, dep2.vulnerabilities, 1)
	vuln = dep2.vulnerabilities[trigger2.Cve]
	assert.Equal(t, trigger2.Cve, vuln.name)
	assert.Equal(t, trigger2.CveLink, vuln.url)
	assert.Equal(t, fmt.Sprintf("%g", trigger2.Cvss2), vuln.cvss2)
	assert.Equal(t, fmt.Sprintf("%g", trigger2.Cvss3), vuln.cvss3)
	assert.Empty(t, dep2.licenses)
}
