package pip

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type PackageMetadata struct {
	Name         string
	Version      string
	Dependencies []string
}

func (j *Job) parsePipList(pipListOutput string) ([]string, error) {
	lines := strings.Split(pipListOutput, "\n")
	packages := []string{}
	for _, line := range lines[2:] {
		fields := strings.Split(line, " ")
		if len(fields) > 0 {
			packages = append(packages, fields[0])
		}
	}
	return packages, nil
}

func (j *Job) parseRequirements() ([]string, error) {
	file, err := os.Open(j.file)
	if err != nil {
		j.err = err
		return nil, err
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

		match := pattern.FindStringSubmatch(line)

		if match != nil {
			packages = append(packages, match[1])
		}
	}

	if err := scanner.Err(); err != nil {
		j.err = err
		return nil, err
	}

	return packages, nil
}

func (j *Job) parseGraph(packages []string, installedPackagesMetadata string) ([]string, []string, []string, error) {
	visitedPackageMetadata := map[string]PackageMetadata{}
	packageMetaData, _ := j.parsePackageMetadata(installedPackagesMetadata)
	missingMetadata := []string{}

	for len(packages) > 0 {
		currentPackage := strings.ToLower(packages[0])
		packages = packages[1:]

		if _, ok := visitedPackageMetadata[currentPackage]; ok {
			continue
		}

		dependencies := packageMetaData[currentPackage].Dependencies
		packages = append(packages, dependencies...)
		if val, ok := packageMetaData[currentPackage]; ok {
			visitedPackageMetadata[currentPackage] = val
		} else {
			missingMetadata = append(missingMetadata, currentPackage)
		}
	}

	nodes := []string{}
	edges := []string{}

	for _, v := range visitedPackageMetadata {
		nodes = append(nodes, fmt.Sprintf("%s %s", v.Name, v.Version))
	}

	fmt.Println("Found", len(visitedPackageMetadata), "installed dependencies.")
	for _, v := range visitedPackageMetadata {
		for _, d := range v.Dependencies {
			edges = append(edges, fmt.Sprintf("%s %s", v.Name, d))
		}
	}

	return nodes, edges, missingMetadata, nil
}

func (j *Job) parsePackageMetadata(installedPackagesMetadata string) (map[string]PackageMetadata, error) {
	result := map[string]PackageMetadata{}
	metadata := strings.Split(installedPackagesMetadata, "---")

	for _, packageMetadata := range metadata {
		lines := strings.Split(packageMetadata, "\n")
		name, version, dependencies := "", "", []string{}

		for _, line := range lines {
			fields := strings.Split(line, ": ")

			if len(fields) == 0 {
				continue
			}

			switch fields[0] {
			case "Name":
				name = fields[1]
			case "Version":
				version = fields[1]
			case "Requires":
				if fields[1] != "" {
					dependencies = strings.Split(fields[1], ", ")
				}
			}
		}

		result[strings.ToLower(name)] = PackageMetadata{name, version, dependencies}
	}
	return result, nil
}