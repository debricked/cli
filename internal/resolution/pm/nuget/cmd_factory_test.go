package nuget

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMakeInstallCmd(t *testing.T) {
	nugetCommand := "dotnet"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeInstallCmd(nugetCommand, "file")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "dotnet")
	assert.Contains(t, args, "restore")
}

func TestMakeInstallCmdPackagsConfig(t *testing.T) {
	nugetCommand := "dotnet"
	cmd, err := CmdFactory{
		execPath: ExecPath{},
	}.MakeInstallCmd(nugetCommand, "testdata/valid/packages.config")
	assert.NoError(t, err)
	assert.NotNil(t, cmd)
	args := cmd.Args
	assert.Contains(t, args, "dotnet")
	assert.Contains(t, args, "restore")

	// Cleanup: Remove the created .csproj file
	if err := os.Remove("testdata/valid/packages.config.csproj"); err != nil {
		t.Fatalf("Failed to remove test file: %v", err)
	}
}
func TestParsePackagesConfig(t *testing.T) {
	tests := []struct {
		name       string
		setup      func() string // function to set up the test environment
		teardown   func()        // function to clean up after the test
		shouldFail bool
	}{
		{
			name: "Non-existent file",
			setup: func() string {
				return "nonexistent_file.config"
			},
			shouldFail: true,
		},
		{
			name: "Unreadable file",
			setup: func() string {
				file, err := os.CreateTemp("", "unreadable_*.config")
				if err != nil {
					t.Fatal(err)
				}
				file.Chmod(0222) // write-only permissions
				return file.Name()
			},
			teardown: func() {
				os.Remove("unreadable_file.config") // clean up the unreadable file
			},
			shouldFail: true,
		},
		{
			name: "Malformed XML",
			setup: func() string {
				file, err := os.CreateTemp("", "malformed_*.config")
				if err != nil {
					t.Fatal(err)
				}
				file.WriteString("malformed xml content")
				return file.Name()
			},
			teardown: func() {
				os.Remove("malformed_file.config") // clean up the malformed file
			},
			shouldFail: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			filePath := test.setup()
			if test.teardown != nil {
				defer test.teardown() // clean up after the test
			}

			_, err := parsePackagesConfig(filePath)
			if (err != nil) != test.shouldFail {
				t.Errorf("parsePackagesConfig() error = %v, shouldFail = %v", err, test.shouldFail)
			}
		})
	}
}

func TestCollectUniqueTargetFrameworks(t *testing.T) {
	packages := []Package{
		{TargetFramework: "net45"},
		{TargetFramework: "net46"},
		{TargetFramework: "net45"},
	}
	got := collectUniqueTargetFrameworks(packages)
	want := "net45;net46"
	if got != want {
		t.Errorf("collectUniqueTargetFrameworks() = %v, want %v", got, want)
	}
}

func TestWriteContentToCsprojFile(t *testing.T) {
	newFilename := "testdata/test_output.csproj"
	content := "<Project></Project>"

	if err := writeContentToCsprojFile(newFilename, content); err != nil {
		t.Fatalf("writeContentToCsprojFile() failed: %v", err)
	}

	if _, err := os.Stat(newFilename); os.IsNotExist(err) {
		t.Fatalf("writeContentToCsprojFile() did not create file")
	}

	// Cleanup: Remove the created file
	if err := os.Remove(newFilename); err != nil {
		t.Fatalf("Failed to remove test file: %v", err)
	}
}

func TestConvertPackagesConfigToCsproj(t *testing.T) {
	tests := []struct {
		filePath  string
		wantError bool
	}{
		{"testdata/valid/packages.config", false},
		{"testdata/invalid/packages.config", true},
	}

	for _, tt := range tests {
		_, err := convertPackagesConfigToCsproj(tt.filePath)
		if (err != nil) != tt.wantError {
			t.Errorf("convertPackagesConfigToCsproj(%q) = %v, want error: %v", tt.filePath, err, tt.wantError)
		}
	}

	// Cleanup: Remove the created .csproj file
	if err := os.Remove("testdata/valid/packages.config.csproj"); err != nil {
		t.Fatalf("Failed to remove test file: %v", err)
	}
}

func TestWriteContentToCsprojFileErr(t *testing.T) {
	tests := []struct {
		name       string
		filename   string
		content    string
		shouldFail bool
		setup      func() // function to set up the environment for the test
		teardown   func() // function to clean up after the test
	}{
		{
			name:       "Valid file name and content",
			filename:   "test.csproj",
			content:    "<Project></Project>",
			shouldFail: false,
			teardown: func() {
				os.Remove("test.csproj") // Clean up the created file
			},
		},
		{
			name:       "Invalid file name",
			filename:   "", // Empty filename is invalid
			content:    "<Project></Project>",
			shouldFail: true,
		},
		{
			name:       "Write to a read-only file",
			filename:   "readonly.csproj",
			content:    "<Project></Project>",
			shouldFail: true,
			setup: func() {
				// Create a read-only file
				file, err := os.Create("readonly.csproj")
				if err != nil {
					panic(err)
				}
				file.Close()
				os.Chmod("readonly.csproj", 0444) // Set file permissions to read-only
			},
			teardown: func() {
				os.Remove("readonly.csproj") // Clean up the read-only file
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.setup != nil {
				test.setup() // set up the environment for the test
			}
			err := writeContentToCsprojFile(test.filename, test.content)
			if (err != nil) != test.shouldFail {
				t.Errorf("writeContentToCsprojFile() error = %v, shouldFail = %v", err, test.shouldFail)
			}
			if test.teardown != nil {
				test.teardown() // clean up after the test
			}
		})
	}
}

func TestCreateCsprojContent(t *testing.T) {
	tests := []struct {
		name                string
		targetFrameworksStr string
		packages            []Package
		shouldFail          bool
		tmpl                string
	}{
		{
			name:                "Valid template action",
			targetFrameworksStr: "netcoreapp3.1",
			packages:            []Package{{ID: "SomePackage", Version: "1.0.0"}},
			shouldFail:          false,
			tmpl:                packagesConfigTemplate,
		},
		{
			name:                "Invalid template action",
			targetFrameworksStr: "netcoreapp3.1",
			packages:            []Package{{ID: "SomePackage", Version: "1.0.0"}},
			shouldFail:          true,
			tmpl: `
		<Project Sdk="Microsoft.NET.Sdk">
			<PropertyGroup>
				<TargetFrameworks>{{.TargetFrameworks}}</TargetFrameworks>
			</PropertyGroup>
			<ItemGroup>
			{{- range .Packages}}
				<PackageReference Include="{{.ID}" Version="{{.Version}}" />  <!-- Missing closing brace -->
			{{- end}}
			</ItemGroup>
		</Project>
		`,
		},
		{
			name:                "Non-existent field",
			targetFrameworksStr: "netcoreapp3.1",
			packages:            []Package{{ID: "SomePackage", Version: "1.0.0"}},
			shouldFail:          true,
			tmpl: `
		<Project Sdk="Microsoft.NET.Sdk">
			<PropertyGroup>
				<TargetFrameworks>{{.NonExistentField}}</TargetFrameworks>  <!-- Non-existent field -->
			</PropertyGroup>
			<ItemGroup>
			{{- range .Packages}}
				<PackageReference Include="{{.ID}}" Version="{{.Version}}" Nofied="{{.Nofield}}"/>   <!-- Non-existent field -->
			{{- end}}
			</ItemGroup>
		</Project>
		`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := createCsprojContentWithTemplate(test.targetFrameworksStr, test.packages, test.tmpl)
			if (err != nil) != test.shouldFail {
				t.Errorf("createCsprojContentWithTemplate() error = %v, shouldFail = %v", err, test.shouldFail)
			}
		})
	}
}
