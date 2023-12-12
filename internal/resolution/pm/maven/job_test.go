package maven

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/maven/testdata"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", CmdFactory{})
	assert.Equal(t, "file", j.GetFile())
	assert.False(t, j.Errors().HasError())
}

func TestRunCmdErr(t *testing.T) {
	cases := []struct {
		error string
		doc   string
	}{
		{
			error: "cmd-error",
			doc:   util.UnknownError,
		},
		{
			error: " |[FATAL] Non-parseable POM /home/asus/Projects/playground/maven-project/pom.xml: end tag name </target> must be the same as start tag <source> from line 37 (position: TEXT seen ...<source>1.6</target>... @37:31)  @ /home/asus/Projects/playground/maven-project/pom.xml, line 37, column 31\n",
			doc:   "Failed to build Maven dependency tree. Your POM file is not valid. Please check /home/asus/Projects/playground/maven-project/pom.xml: end tag name </target> must be the same as start tag <source> from line 37 (position: TEXT seen ...<source>1.6</target>... @37:31)  @ /home/asus/Projects/playground/maven-project/pom.xml, line 37, column 31",
		},
		{
			error: "  |[WARNING] Failed to retrieve plugin descriptor for org.apache.maven.plugins:maven-compiler-plugin:2.3.2: Plugin org.apache.maven.plugins:maven-compiler-plugin:2.3.2 or one of its dependencies could not be resolved: Failed to read artifact descriptor for org.apache.maven.plugins:maven-compiler-plugin:jar:2.3.2\n",
			doc:   "We weren't able to retrieve one or more plugin descriptor(s). Please check your Internet connection and try again.",
		},
		{
			error: "   |[ERROR] 'dependencies.dependency.version' for org.hamcrest:hamcrest-library:jar must not contain any of these characters \\/:\"<>|?* but found * @ com.example.maven-project:maven-project:1.0-SNAPSHOT, /home/asus/Projects/playground/maven-project/pom.xml, line 196, column 18\n",
			doc:   "There is an error in dependencies: 'dependencies.dependency.version' for org.hamcrest:hamcrest-library:jar must not contain any of these characters \\/:\"<>|?* but found *",
		},
		{
			error: "    |[ERROR] Failed to execute goal on project jackpot: Could not resolve dependencies for project com.jeteo:jackpot:war:1.0-SNAPSHOT: The following artifacts could not be resolved: javax.servlet:com.springsource.javax.servlet:jar:2.5.0, javax.servlet:com.springsource.javax.servlet.jsp.jstl:jar:1.2.0 (http://repository.springsource.com/maven/bundles/release) -> [Help 1]\n",
			doc:   "Could not resolve dependencies for project com.jeteo:jackpot:war:1.0-SNAPSHOT: The following artifacts could not be resolved: javax.servlet:com.springsource.javax.servlet:jar:2.5.0, javax.servlet:com.springsource.javax.servlet.jsp.jstl:jar:1.2.0  \nTry to run `mvn dependency:tree -e` to get more details. If this is a private dependency, make sure you have access to install it.",
		},
	}

	for _, c := range cases {
		expectedError := util.NewPMJobError(c.error)
		expectedError.SetDocumentation(c.doc)

		cmdErr := errors.New(c.error)
		j := NewJob("file", testdata.CmdFactoryMock{Err: cmdErr})

		go jobTestdata.WaitStatus(j)

		j.Run()

		assert.Len(t, j.Errors().GetAll(), 1)
		assert.Contains(t, j.Errors().GetAll(), expectedError)
	}
}

func TestRunCmdOutputErr(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "bad-name"})

	go jobTestdata.WaitStatus(j)

	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunCmdOutputErrNoOutput(t *testing.T) {
	j := NewJob("file", testdata.CmdFactoryMock{Name: "go", Arg: "bad-arg"})

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
	j := NewJob("file", testdata.CmdFactoryMock{Name: "echo"})

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.False(t, j.Errors().HasError())
}
