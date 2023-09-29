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
	"sort"
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

type CmdFactory struct {
	execPath               IExecPath
	packageConfgRegex      string
	packagesConfigTemplate string
}

func NewCmdFactory(execPath IExecPath) CmdFactory {
	return CmdFactory{
		execPath:               execPath,
		packageConfgRegex:      PackagesConfigRegex,
		packagesConfigTemplate: packagesConfigTemplate,
	}
}

func (cmdf CmdFactory) MakeInstallCmd(command string, file string) (*exec.Cmd, error) {

	path, err := cmdf.execPath.LookPath(command)

	if err != nil {
		return nil, err
	}

	// If the file is a packages.config file, convert it to a .csproj file
	// check regex with PackagesConfigRegex
	packageConfig, err := regexp.Compile(cmdf.packageConfgRegex)
	if err != nil {
		return nil, err
	}

	if packageConfig.MatchString(file) {
		file, err = cmdf.convertPackagesConfigToCsproj(file, command)
		if err != nil {
			return nil, err
		}
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
func (cmdf CmdFactory) convertPackagesConfigToCsproj(filePath string, command string) (string, error) {
	packages, err := parsePackagesConfig(filePath)
	if err != nil {
		return "", err
	}

	targetFrameworksStr, err := collectUniqueTargetFrameworks(packages.Packages, command)
	if err != nil {
		return "", err
	}
	csprojContent, err := cmdf.createCsprojContentWithTemplate(targetFrameworksStr, packages.Packages)
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

var ioReadAllCsproj = io.ReadAll

func getDotnetVersion(command string) (string, error) {
	cmd := exec.Command(command, "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func getDefaultFrameworkOfDotnetVersion(dotnetVersion string) string {

	if strings.HasPrefix(dotnetVersion, "7") {
		return "net7.0"
	} else if strings.HasPrefix(dotnetVersion, "6") {
		return "net6.0"
	} else if strings.HasPrefix(dotnetVersion, "5") {
		return "net5.0"
	} else if strings.HasPrefix(dotnetVersion, "3") {
		return "netcoreapp3.1"
	} else if strings.HasPrefix(dotnetVersion, "2") {
		return "netcoreapp2.1"
	} else if strings.HasPrefix(dotnetVersion, "1") {
		return "netcoreapp1.1"
	}

	return "net6.0"
}

func parsePackagesConfig(filePath string) (*Packages, error) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	byteValue, err := ioReadAllCsproj(xmlFile)
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

func collectUniqueTargetFrameworks(packages []Package, command string) (string, error) {
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

	sort.Strings(targetFrameworks) // Sort the targetFrameworks slice

	if len(targetFrameworks) == 0 {
		dotnetVersion, err := getDotnetVersion(command)
		if err != nil {
			return "", err
		}

		targetFrameworks = append(targetFrameworks, getDefaultFrameworkOfDotnetVersion(dotnetVersion))
	}

	return strings.Join(targetFrameworks, ";"), nil
}
func (cmdf CmdFactory) createCsprojContentWithTemplate(targetFrameworksStr string, packages []Package) (string, error) {
	tmplParsed, err := template.New("csproj").Parse(cmdf.packagesConfigTemplate)
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

var osCreateCsproj = os.Create

func writeContentToCsprojFile(newFilename string, content string) error {

	csprojFile, err := osCreateCsproj(newFilename)
	if err != nil {
		return err
	}
	defer csprojFile.Close()

	_, err = csprojFile.WriteString(content)

	return err
}
