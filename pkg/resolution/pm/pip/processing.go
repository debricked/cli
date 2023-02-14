package pip

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type PackageMetadata struct {
	name         string
	version      string
	dependencies []string
}

func (j *Job) parseRequirements() ([]string, error) {
	file, err := os.Open(j.file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	packages := []string{}

	pattern := regexp.MustCompile(`^([^\s]+?)(?:[=<>!~]+(.+))?$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Match package name and version using regular expression
		match := pattern.FindStringSubmatch(line)
		if match != nil {
			packages = append(packages, match[1])
		}
	}

	// Check for any scanning errors
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return packages, nil
}

func (j *Job) parseGraph(packages []string, installedPackagesMetadata string) ([]string, []string, error) {

	visitedPackagesVersions := map[string]string{}
	visitedPackagesRelations := map[string][]string{}

	PackageMetadata, _ := j.parsePackageMetadata(installedPackagesMetadata)
	fmt.Println("metadata", PackageMetadata)
	fmt.Println("metadata", PackageMetadata)

	for len(packages) > 0 {
		if _, ok := visitedPackagesVersions[packages[0]]; ok {
			packages = packages[1:]
			continue
		}
		p := packages[0]
		packages = packages[1:]

		pm := PackageMetadata[p]
		fmt.Println("pm", p, pm)

		version := PackageMetadata[p].version
		dependencies := PackageMetadata[p].dependencies

		packages = append(packages, dependencies...)
		visitedPackagesVersions[p] = version
		visitedPackagesRelations[p] = dependencies

	}
	//transform maps to list of strings
	nodes := []string{}
	edges := []string{}

	for k, v := range visitedPackagesVersions {
		nodes = append(nodes, fmt.Sprintf("%s %s", k, v))
	}

	fmt.Println("Visited", visitedPackagesRelations)
	for k, v := range visitedPackagesRelations {
		for _, d := range v {
			edges = append(edges, fmt.Sprintf("%s %s", k, d))
		}
	}

	// return two bytes arrays

	return nodes, edges, nil
}

func (j *Job) parsePackageMetadata(installedPackagesMetadata string) (map[string]PackageMetadata, error) {

	m := map[string]PackageMetadata{}

	metadata := strings.Split(installedPackagesMetadata, "---")
	for _, packageMetadata := range metadata {
		lines := strings.Split(packageMetadata, "\n")

		name, version, dependencies := "", "", []string{}
		for _, line := range lines {
			fields := strings.Split(line, ":")
			if len(fields) == 0 {
				continue
			}
			if fields[0] == "Name" {
				name = fields[1]
			}
			if fields[0] == "Version" {
				version = fields[1]
			}
			if fields[0] == "Requires" {
				dependencies = strings.Split(fields[1], ",")
			}
		}

		m[name] = PackageMetadata{name, version, dependencies}
	}

	return m, nil
}

func (j *Job) parsePipList(pipListOutput string) ([]string, error) {

	lines := strings.Split(pipListOutput, "\n")
	packages := []string{}
	//skip 2 lines
	for _, line := range lines[2:] {
		fields := strings.Split(line, " ")
		if len(fields) > 0 {
			packages = append(packages, fields[0])
		}
	}
	return packages, nil
}
