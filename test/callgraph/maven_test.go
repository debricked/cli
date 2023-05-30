package resolve

// import (
// 	"os"
// 	"path"
// 	"testing"

// 	"github.com/debricked/cli/internal/callgraph/cgexec"
// 	// "github.com/debricked/cli/internal/callgraph/language"
// 	"github.com/stretchr/testify/assert"
// )

// func TestMakeBuildMavenCmdFunctional(t *testing.T) {
// 	workingDir, err := os.Getwd()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	t.Log(workingDir)
// 	javaProjectPath := "testdata/mvnproj"
// 	javaProjectAbsPath := path.Join(workingDir, javaProjectPath)
// 	javaProjectTargetAbsPath := path.Join(javaProjectAbsPath, "target")
// 	assert.NoDirExists(t, javaProjectTargetAbsPath)
// 	ctx, _ := cgexec.NewContext(10000)
// 	cmd, err := CmdFactory{}.MakeBuildMavenCmd(javaProjectAbsPath, ctx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = cgexec.RunCommand(cmd, ctx)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	assert.DirExists(t, javaProjectTargetAbsPath)
// 	os.RemoveAll(javaProjectTargetAbsPath)
// 	assert.NoDirExists(t, javaProjectTargetAbsPath)
// }

// func TestResolveMaven(t *testing.T) {
// 	cases := []struct {
// 		name             string
// 		requirementsFile string
// 		expectedFile     string
// 	}{
// 		{
// 			name:             "basic requirements.txt",
// 			requirementsFile: "testdata/pip/requirements.txt",
// 			expectedFile:     "testdata/pip/expected.lock",
// 		},
// 	}

// 	for _, c := range cases {
// 		t.Run(c.name, func(t *testing.T) {
// 			resolveCmd := resolve.NewResolveCmd(wire.GetCliContainer().Resolver())
// 			err := resolveCmd.RunE(resolveCmd, []string{c.requirementsFile})
// 			assert.NoError(t, err)

// 			lockFileDir := filepath.Dir(c.requirementsFile)
// 			lockFile := filepath.Join(lockFileDir, ".requirements.txt.debricked.lock")
// 			lockFileContents, fileErr := os.ReadFile(lockFile)
// 			assert.NoError(t, fileErr)

// 			expectedFileContents, fileErr := os.ReadFile(c.expectedFile)
// 			assert.NoError(t, fileErr)

// 			assert.Equal(t, string(expectedFileContents), string(lockFileContents))
// 		})
// 	}
// }
