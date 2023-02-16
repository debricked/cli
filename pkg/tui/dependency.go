package tui

import (
	"fmt"
	"github.com/debricked/cli/pkg/automation"
)

type dependency struct {
	name            string
	url             string
	vulnerabilities map[string]vulnerability
	licenses        map[string]bool
}

func makeDependenciesFromTriggers(triggers []automation.TriggerEvent) map[string]dependency {
	dependencies := map[string]dependency{}

	for _, trigger := range triggers {
		dep, ok := dependencies[trigger.Dependency]
		if !ok {
			dep = dependency{
				name:            trigger.Dependency,
				url:             trigger.DependencyLink,
				vulnerabilities: map[string]vulnerability{},
				licenses:        map[string]bool{},
			}
			dependencies[dep.name] = dep
		}

		for _, license := range trigger.Licenses {
			if len(license) > 0 {
				dep.licenses[license] = true
			}
		}

		if _, ok = dep.vulnerabilities[trigger.Cve]; !ok && len(trigger.Cve) > 0 {
			dep.vulnerabilities[trigger.Cve] = vulnerability{
				name:  trigger.Cve,
				url:   trigger.CveLink,
				cvss2: fmt.Sprintf("%g", trigger.Cvss2),
				cvss3: fmt.Sprintf("%g", trigger.Cvss3),
			}
		}
	}

	return dependencies
}
