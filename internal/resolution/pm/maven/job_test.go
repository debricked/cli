package maven

import (
	"errors"
	"path/filepath"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/maven/testdata"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", CmdFactory{}, PomService{})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunCmdErr(t *testing.T) {
	cases := []struct {
		name  string
		error string
		doc   string
	}{
		{
			name:  "General error",
			error: "cmd-error",
			doc:   util.UnknownError,
		},

		{
			name:  "Mvn not found",
			error: "        |exec: \"mvn\": executable file not found in $PATH",
			doc:   "Mvn wasn't found. Please check if it is installed and accessible by the CLI.",
		},
		{
			name:  "Invalid XML",
			error: " |[FATAL] Non-parseable POM /home/asus/Projects/playground/maven-project/pom.xml: end tag name </target> must be the same as start tag <source> from line 37 (position: TEXT seen ...<source>1.6</target>... @37:31)  @ /home/asus/Projects/playground/maven-project/pom.xml, line 37, column 31\n",
			doc:   "Failed to build Maven dependency tree. Your POM file is not valid. Please check /home/asus/Projects/playground/maven-project/pom.xml: end tag name </target> must be the same as start tag <source> from line 37 (position: TEXT seen ...<source>1.6</target>... @37:31)  @ /home/asus/Projects/playground/maven-project/pom.xml, line 37, column 31",
		},
		{
			name:  "No Internet",
			error: "  |[WARNING] Failed to retrieve plugin descriptor for org.apache.maven.plugins:maven-compiler-plugin:2.3.2: Plugin org.apache.maven.plugins:maven-compiler-plugin:2.3.2 or one of its dependencies could not be resolved: Failed to read artifact descriptor for org.apache.maven.plugins:maven-compiler-plugin:jar:2.3.2\n",
			doc:   "We weren't able to retrieve one or more plugin descriptor(s). Please check your Internet connection and try again.",
		},
		{
			name:  "Invalid dependency",
			error: "   |[ERROR] 'dependencies.dependency.version' for org.hamcrest:hamcrest-library:jar must not contain any of these characters \\/:\"<>|?* but found * @ com.example.maven-project:maven-project:1.0-SNAPSHOT, /home/asus/Projects/playground/maven-project/pom.xml, line 196, column 18\n",
			doc:   "There is an error in dependencies: 'dependencies.dependency.version' for org.hamcrest:hamcrest-library:jar must not contain any of these characters \\/:\"<>|?* but found *",
		},
		{
			name:  "Invalid version",
			error: "    |[ERROR] Failed to execute goal on project jackpot: Could not resolve dependencies for project com.jeteo:jackpot:war:1.0-SNAPSHOT: The following artifacts could not be resolved: javax.servlet:com.springsource.javax.servlet:jar:2.5.0, javax.servlet:com.springsource.javax.servlet.jsp.jstl:jar:1.2.0 (http://repository.springsource.com/maven/bundles/release) -> [Help 1]\n",
			doc:   "An error occurred during dependencies resolve for: com.jeteo:jackpot:war:1.0-SNAPSHOT\nTry to run `mvn dependency:tree -e` to get more details.\nIf this is a private dependency, please make sure that the debricked CLI has access to install it or pre-install it before running the debricked CLI.",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expectedError := util.NewPMJobError(c.error)
			expectedError.SetDocumentation(c.doc)

			cmdErr := errors.New(c.error)
			j := NewJob("file", testdata.CmdFactoryMock{Err: cmdErr}, testdata.PomServiceMock{})

			go jobTestdata.WaitStatus(j)

			j.Run()

			allErrors := j.Errors().GetAll()

			assert.Len(t, allErrors, 1)
			assert.Contains(t, allErrors, expectedError)
		})
	}
}

func TestRunCmdOutputErr(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "bad-name"}, testdata.PomServiceMock{})

	go jobTestdata.WaitStatus(j)

	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunCmdOutputErrNoOutput(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "go", Arg: "bad-arg"}, testdata.PomServiceMock{})

	go jobTestdata.WaitStatus(j)

	j.Run()

	errs := j.Errors().GetAll()
	assert.Len(t, errs, 1)
	err := errs[0]

	// assert empty because, when Output is executed it will allocate memory for the byte slice to contain the standard output.
	// However since no bytes are sent to standard output err will be empty here.
	assert.Empty(t, err.Error())
}

func TestRun(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "echo"}, testdata.PomServiceMock{})

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.False(t, j.Errors().HasError())
}

func TestSuccessfulRunWithActualFile(t *testing.T) {
	cases := []struct {
		name string
		file string
	}{
		{
			name: "valid file",
			file: filepath.Join("testdata", "pom.xml"),
		},
		{
			name: "valid child",
			file: filepath.Join("testdata", "guava", "pom.xml"),
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			j := NewJob(c.file, testdata.CmdFactoryMock{Name: "echo"}, PomService{})

			go jobTestdata.WaitStatus(j)

			j.Run()

			assert.False(t, j.Errors().HasError())
		})
	}
}

func TestRunWithActualFileErrOutput(t *testing.T) {
	cases := []struct {
		name  string
		file  string
		error string
		doc   string
	}{
		{
			name:  "not a pom",
			file:  filepath.Join("testdata", "notAPom.xml"),
			error: "EOF",
			doc:   "This file doesn't contain valid XML",
		},
		{
			name:  "invalid pom",
			file:  filepath.Join("testdata", "invalidPom.xml"),
			error: "XML syntax error on line 13: invalid characters between </artifactId and >",
			doc:   "XML syntax error on line 13: invalid characters between </artifactId and >",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			j := NewJob(c.file, testdata.CmdFactoryMock{Name: "echo"}, PomService{})

			go jobTestdata.WaitStatus(j)

			j.Run()

			allErrors := j.Errors().GetAll()

			expectedError := util.NewPMJobError(c.error)
			expectedError.SetStatus("parsing XML")
			expectedError.SetDocumentation(c.doc)

			assert.Len(t, allErrors, 1)
			assert.Contains(t, allErrors, expectedError)
		})
	}
}
