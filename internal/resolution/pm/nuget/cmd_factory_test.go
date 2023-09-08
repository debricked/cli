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
		filePath  string
		wantError bool
	}{
		{"testdata/valid/packages.config", false},
		{"testdata/invalid/packages.config", true},
	}

	for _, tt := range tests {
		_, err := parsePackagesConfig(tt.filePath)
		if (err != nil) != tt.wantError {
			t.Errorf("parsePackagesConfig(%q) = %v, want error: %v", tt.filePath, err, tt.wantError)
		}
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

func TestCreateCsprojContent(t *testing.T) {
	packages := []Package{
		{ID: "Test.Package.1", Version: "1.0.0", TargetFramework: "net45"},
		{ID: "Test.Package.2", Version: "2.0.0", TargetFramework: "net46"},
	}
	targetFrameworksStr := "net45;net46"

	got, err := createCsprojContent(targetFrameworksStr, packages)
	if err != nil {
		t.Fatalf("createCsprojContent() failed: %v", err)
	}

	// We're just checking if the function returns a non-empty string
	// For a more rigorous test, we'd check the exact content of the string
	if got == "" {
		t.Errorf("createCsprojContent() returned an empty string")
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
