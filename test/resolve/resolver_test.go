package resolve

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/debricked/cli/internal/cmd/resolve"
	"github.com/debricked/cli/internal/resolution/pm/npm"
	"github.com/debricked/cli/internal/wire"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestResolves(t *testing.T) {
	cases := []struct {
		name           string
		manifestFile   string
		lockFileName   string
		packageManager string
	}{
		{
			name:           "basic composer.json",
			manifestFile:   "testdata/composer/composer.json",
			lockFileName:   "composer.lock",
			packageManager: "composer",
		},
		{
			name:           "basic package.json (Yarn)",
			manifestFile:   "testdata/npm/package.json",
			lockFileName:   "yarn.lock",
			packageManager: "yarn",
		},
		{
			name:           "basic package.json (NPM)",
			manifestFile:   "testdata/npm/package.json",
			lockFileName:   "package-lock.json",
			packageManager: "npm",
		},
		{
			name:           "basic bower.json",
			manifestFile:   "testdata/bower/bower.json",
			lockFileName:   "bower.debricked.lock",
			packageManager: "bower",
		},
		{
			name:           "basic requirements.txt",
			manifestFile:   "testdata/pip/requirements.txt",
			lockFileName:   "requirements.txt.pip.debricked.lock",
			packageManager: "pip",
		},
		{
			name:           "basic .csproj",
			manifestFile:   "testdata/nuget/csproj/basic.csproj",
			lockFileName:   "packages.lock.json",
			packageManager: "nuget",
		},
		{
			name:           "basic packages.config",
			manifestFile:   "testdata/nuget/packagesconfig/packages.config",
			lockFileName:   "packages.config.nuget.debricked.lock",
			packageManager: "nuget",
		},
		{
			name:           "basic go.mod",
			manifestFile:   "testdata/gomod/go.mod",
			lockFileName:   "gomod.debricked.lock",
			packageManager: "gomod",
		},
		{
			name:           "basic pom.xml",
			manifestFile:   "testdata/maven/pom.xml",
			lockFileName:   "maven.debricked.lock",
			packageManager: "maven",
		},
		{
			name:           "basic build.gradle",
			manifestFile:   "testdata/gradle/build.gradle",
			lockFileName:   "gradle.debricked.lock",
			packageManager: "gradle",
		},
	}

	for _, cT := range cases {
		c := cT
		t.Run(c.name, func(t *testing.T) {
			if c.packageManager == npm.Name {
				viper.Set(resolve.NpmPreferredFlag, true)
			}

			resolveCmd := resolve.NewResolveCmd(wire.GetCliContainer().Resolver())
			lockFileDir := filepath.Dir(c.manifestFile)
			lockFile := filepath.Join(lockFileDir, c.lockFileName)
			// Remove the lock file if it exists
			os.Remove(lockFile)

			err := resolveCmd.RunE(resolveCmd, []string{c.manifestFile})
			assert.NoError(t, err)

			lockFileContents, fileErr := os.ReadFile(lockFile)
			assert.NoError(t, fileErr)

			actualString := string(lockFileContents)

			assert.Greater(t, len(actualString), 0)

		})
	}
}

// TestResolvePub verifies that `dart pub get` produces a pubspec.lock when
// the Dart SDK is available on PATH. The test is skipped if dart is not found,
// so it is safe to run in environments without the Dart SDK installed.
func TestResolvePub(t *testing.T) {
	if _, err := exec.LookPath("dart"); err != nil {
		t.Skip("dart not found in PATH; skipping pub resolution test")
	}

	manifestFile := "testdata/pub/pubspec.yaml"
	lockFile := "testdata/pub/pubspec.lock"
	depsFile := "testdata/pub/pubspec.deps.json"

	// Preserve and restore the original lock file if it exists so the test
	// data directory remains clean after the test run.
	original, readErr := os.ReadFile(lockFile)
	originalDeps, readDepsErr := os.ReadFile(depsFile)

	t.Cleanup(func() {
		if readErr == nil {
			_ = os.WriteFile(lockFile, original, 0600)
		} else {
			_ = os.Remove(lockFile)
		}

		if readDepsErr == nil {
			_ = os.WriteFile(depsFile, originalDeps, 0600)
		} else {
			_ = os.Remove(depsFile)
		}
	})

	// Remove any stale lock file so the resolver has to generate a fresh one.
	_ = os.Remove(lockFile)
	_ = os.Remove(depsFile)

	resolveCmd := resolve.NewResolveCmd(wire.GetCliContainer().Resolver())
	err := resolveCmd.RunE(resolveCmd, []string{manifestFile})
	assert.NoError(t, err)

	contents, fileErr := os.ReadFile(lockFile)
	assert.NoError(t, fileErr)
	assert.Greater(t, len(contents), 0)

	depsContents, depsErr := os.ReadFile(depsFile)
	assert.NoError(t, depsErr)
	assert.Greater(t, len(depsContents), 0)
}
