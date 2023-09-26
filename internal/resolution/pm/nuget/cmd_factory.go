package nuget

import (
	"bytes"
	"encoding/xml"
	"html/template"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

type ICmdFactory interface {
	MakeInstallCmd(command string, file string) (*exec.Cmd, error)
}

type IExecPath interface {
	LookPath(file string) (string, error)
}

type ExecPath struct {
}

func (ExecPath) LookPath(file string) (string, error) {
	return exec.LookPath(file)
}

type CmdFactory struct {
	execPath IExecPath
}

var packagesConfigTemplate = `
<Project Sdk="Microsoft.NET.Sdk">
	<PropertyGroup>
		<TargetFrameworks>{{.TargetFrameworks}}</TargetFrameworks>
	</PropertyGroup>
	<ItemGroup>
	{{- range .Packages}}
		<PackageReference Include="{{.ID}}" Version="{{.Version}}" />
	{{- end}}
	</ItemGroup>
</Project>
`

func (cmdf CmdFactory) MakeInstallCmd(command string, file string) (*exec.Cmd, error) {

	// If the file is a packages.config file, convert it to a .csproj file
	// check regex with PackagesConfigRegex
	packageConfig, err := regexp.Compile(PackagesConfigRegex)
	if err != nil {
		return nil, err
	}

	if packageConfig.MatchString(file) {
		file, err = convertPackagesConfigToCsproj(file)
		if err != nil {
			return nil, err
		}
	}

	path, err := cmdf.execPath.LookPath(command)

	if err != nil {
		return nil, err
	}

	fileDir := filepath.Dir(file)

	return &exec.Cmd{
		Path: path,
		Args: []string{command, "restore",
			"--use-lock-file",
		},
		Dir: fileDir,
	}, err
}

type Packages struct {
	Packages []Package `xml:"package"`
}

type Package struct {
	ID              string `xml:"id,attr"`
	Version         string `xml:"version,attr"`
	TargetFramework string `xml:"targetFramework,attr"`
}

// convertPackagesConfigtoCsproj converts a packages.config file to a .csproj file
// this is to enable the use of the dotnet restore command
// that enables debricked to parse out transitive dependencies.
// This may add some additional framework dependencies that will not show up if
// we only scan the packages.config file.
func convertPackagesConfigToCsproj(filePath string) (string, error) {
	packages, err := parsePackagesConfig(filePath)
	if err != nil {
		return "", err
	}

	targetFrameworksStr := collectUniqueTargetFrameworks(packages.Packages)
	csprojContent, err := createCsprojContentWithTemplate(targetFrameworksStr, packages.Packages, packagesConfigTemplate)
	if err != nil {
		return "", err
	}

	newFilename := filePath + ".csproj"
	err = writeContentToCsprojFile(newFilename, csprojContent)
	if err != nil {
		return "", err
	}

	return newFilename, nil
}

func parsePackagesConfig(filePath string) (*Packages, error) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}

	var packages Packages
	err = xml.Unmarshal(byteValue, &packages)
	if err != nil {
		return nil, err
	}

	return &packages, nil
}

func collectUniqueTargetFrameworks(packages []Package) string {
	uniqueTargetFrameworks := make(map[string]struct{})
	for _, pkg := range packages {
		uniqueTargetFrameworks[pkg.TargetFramework] = struct{}{}
	}

	var targetFrameworks []string
	for framework := range uniqueTargetFrameworks {
		if framework != "" {
			targetFrameworks = append(targetFrameworks, framework)
		}
	}

	return strings.Join(targetFrameworks, ";")
}

func createCsprojContentWithTemplate(targetFrameworksStr string, packages []Package, tmpl string) (string, error) {
	tmplParsed, err := template.New("csproj").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	err = tmplParsed.Execute(&tpl, map[string]interface{}{
		"TargetFrameworks": targetFrameworksStr,
		"Packages":         packages,
	})
	if err != nil {
		return "", err
	}

	return tpl.String(), nil
}

func writeContentToCsprojFile(newFilename string, content string) error {

	csprojFile, err := os.Create(newFilename)
	if err != nil {
		return err
	}
	defer csprojFile.Close()

	_, err = csprojFile.WriteString(content)

	return err
}
