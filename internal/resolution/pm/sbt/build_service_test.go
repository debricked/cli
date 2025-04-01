package sbt

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseBuildModules(t *testing.T) {
	content := `
		name := "test-project"
		version := "1.0.0"
		
		lazy val core = project("core")
		  .settings(
		    libraryDependencies += "org.scala-lang" % "scala-library" % "2.13.8"
		  )
		
		lazy val api = project("api")
		  .dependsOn(core)
		  .settings(
		    libraryDependencies += "com.typesafe.akka" %% "akka-http" % "10.2.9"
		  )
		
		lazy val root = (project in file("."))
		  .aggregate(core, api)
	`

	tmpFile, err := os.CreateTemp("", "build.sbt")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	b := BuildService{}
	modules, err := b.ParseBuildModules(tmpFile.Name())

	assert.Nil(t, err)
	assert.Contains(t, modules, "core")
	assert.Contains(t, modules, "api")
}

func TestParseBuildModulesInvalidFile(t *testing.T) {
	b := BuildService{}
	_, err := b.ParseBuildModules("non_existent_file.sbt")

	assert.NotNil(t, err)
}

func TestFindPomFile(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sbt-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	scalaDir := filepath.Join(tempDir, "target", "scala-2.13")
	err = os.MkdirAll(scalaDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directories: %v", err)
	}

	pomPath := filepath.Join(scalaDir, "test-project-1.0.0.pom")
	err = os.WriteFile(pomPath, []byte("<project></project>"), 0600)
	if err != nil {
		t.Fatalf("Failed to create pom file: %v", err)
	}

	foundPom, err := FindPomFile(tempDir)

	assert.Nil(t, err)
	assert.Equal(t, pomPath, foundPom)
}

func TestFindPomFileNoTarget(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sbt-test-no-target")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	foundPom, err := FindPomFile(tempDir)

	assert.Nil(t, err)
	assert.Empty(t, foundPom)
}

func TestRenamePomToXml(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "sbt-rename-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	pomContent := "<project><artifactId>test</artifactId></project>"
	pomPath := filepath.Join(tempDir, "test.pom")
	err = os.WriteFile(pomPath, []byte(pomContent), 0600)
	if err != nil {
		t.Fatalf("Failed to create pom file: %v", err)
	}

	xmlPath, err := RenamePomToXml(pomPath, tempDir)

	assert.Nil(t, err)
	assert.Equal(t, filepath.Join(tempDir, "pom.xml"), xmlPath)

	content, err := os.ReadFile(xmlPath)
	assert.Nil(t, err)
	assert.Equal(t, pomContent, string(content))
}
