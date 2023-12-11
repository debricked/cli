package gradle

import (
	"errors"
	"testing"

	jobTestdata "github.com/debricked/cli/internal/resolution/job/testdata"
	"github.com/debricked/cli/internal/resolution/pm/gradle/testdata"
	"github.com/debricked/cli/internal/resolution/pm/util"
	"github.com/debricked/cli/internal/resolution/pm/writer"
	writerTestdata "github.com/debricked/cli/internal/resolution/pm/writer/testdata"
	"github.com/stretchr/testify/assert"
)

func TestNewJob(t *testing.T) {
	j := NewJob("file", "dir", "nil", "nil", CmdFactory{}, writer.FileWriter{})
	assert.Equal(t, "file", j.GetFile())
	assert.Equal(t, "dir", j.GetDir())
	assert.False(t, j.Errors().HasError())
}

func TestRunCmdErr(t *testing.T) {
	cases := []struct {
		error string
		doc   string
	}{
		{
			error: "cmd-error",
			doc:   "No specific documentation for this problem yet, please report it to us! :)",
		},
		{
			error: "* What went wrong:\nCould not open init remapped class cache for 60sdrkd1iuvns7c8vzs3hv858 (/home/asus/.gradle/caches/5.4/scripts-remapped/_gradle_init_script_debricked_9ll3l6asw7d59x4iljlnzgcpd/60sdrkd1iuvns7c8vzs3hv858/inita22655f7e805aaeb10a177dc56aa75ac).\n> Could not open init generic class cache for initialization script '/home/asus/Projects/playground/gradle-retrolambda/.gradle-init-script.debricked.groovy' (/home/asus/.gradle/caches/5.4/scripts/60sdrkd1iuvns7c8vzs3hv858/init/inita22655f7e805aaeb10a177dc56aa75ac).\n   > BUG! exception in phase 'semantic analysis' in source unit '_BuildScript_' Unsupported class file major version 57\n",
			doc:   "Failed to build Gradle dependency tree. The process has failed with following error: exception in phase 'semantic analysis' in source unit '_BuildScript_' Unsupported class file major version 57. Try run following command to get stacktrace: `nil --init-script nil debrickedFindSubProjectPaths --stacktrace` Replace --stacktrace with --info or --debug option to get more log output. Or with --scan to get full insights.",
		},
		{
			error: "  |Error: Could not find or load main class org.gradle.wrapper.GradleWrapperMain\n        |Caused by: java.lang.ClassNotFoundException: org.gradle.wrapper.GradleWrapperMain\n",
			doc:   "Failed to build Gradle dependency tree. The process has failed with following error: Could not find or load main class org.gradle.wrapper.GradleWrapperMain. You probably trying to run the command not from the root directory.",
		},
		{
			error: "  |* What went wrong:\n        |Project directory '/home/asus/Projects/playground/protobuf-gradle-plugin/testProjectLite' is not part of the build defined by settings file '/home/asus/Projects/playground/protobuf-gradle-plugin/settings.gradle'. If this is an unrelated build, it must have its own settings file.",
			doc:   "Failed to build Gradle dependency tree. The process has failed with following error: Project directory '/home/asus/Projects/playground/protobuf-gradle-plugin/testProjectLite' is not part of the build defined by settings file '/home/asus/Projects/playground/protobuf-gradle-plugin/settings.gradle'. This error might be caused by inclusion of test folders into resolve process. Try running resolve command with -e flag. For example, `debricked resolve -e \"**/test*/**\"` will exclude all folders that start from 'test' from resolution process. Or if this is an unrelated build, it must have its own settings file.",
		},
		{
			error: "  |A problem occurred evaluating settings 'protobuf-gradle-plugin'.\n        |> Could not get unknown property 'glkjhe' for settings 'protobuf-gradle-plugin' of type org.gradle.initialization.DefaultSettings.",
			doc:   "Failed to build Gradle dependency tree. The process has failed with following error: Could not get unknown property 'glkjhe' for settings 'protobuf-gradle-plugin' of type org.gradle.initialization.DefaultSettings.. Please check your settings.gradle file for errors.",
		},
	}

	for _, c := range cases {
		expectedError := util.NewPMJobError(c.error)
		expectedError.SetDocumentation(c.doc)

		cmdErr := errors.New(c.error)
		j := NewJob("file", "dir", "nil", "nil", testdata.CmdFactoryMock{Err: cmdErr}, writer.FileWriter{})

		go jobTestdata.WaitStatus(j)

		j.Run()

		assert.Len(t, j.Errors().GetAll(), 1)
		assert.Contains(t, j.Errors().GetAll(), expectedError)
	}
}

func TestRunCmdOutputErr(t *testing.T) {
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: errors.New("create-error")}

	j := NewJob("file", "dir", "gradlew", "path", testdata.CmdFactoryMock{Name: "bad-name"}, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	jobTestdata.AssertPathErr(t, j.Errors())
}

func TestRunCreateErr(t *testing.T) {
	createErr := errors.New("create-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CreateErr: createErr}
	j := NewJob("file", "dir", "gradlew", "path", testdata.CmdFactoryMock{Name: "echo", Err: createErr}, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), util.NewPMJobError(createErr.Error()))
}

func TestRunWriteErr(t *testing.T) {
	writeErr := errors.New("write-error")
	fileWriterMock := &writerTestdata.FileWriterMock{WriteErr: writeErr}
	j := NewJob("file", "dir", "", "", testdata.CmdFactoryMock{Name: "echo", Err: writeErr}, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), util.NewPMJobError(writeErr.Error()))
}

func TestRunCloseErr(t *testing.T) {
	closeErr := errors.New("close-error")
	fileWriterMock := &writerTestdata.FileWriterMock{CloseErr: closeErr}
	j := NewJob("file", "dir", "gradlew", "path", testdata.CmdFactoryMock{Name: "echo", Err: closeErr}, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
	assert.Contains(t, j.Errors().GetAll(), util.NewPMJobError(closeErr.Error()))
}

func TestRunPermissionFailBeforeOutputErr(t *testing.T) {
	permissionErr := errors.New("give-error-on-gradle gradlew\": permission denied")
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", "dir", "gradlew", "path", testdata.CmdFactoryMock{Name: "echo", Err: permissionErr}, fileWriterMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 2)
}

func TestRunPermissionErr(t *testing.T) {
	permissionErr := errors.New("asdhjaskdhqwe gradlew\": permission denied")
	fileWriterMock := &writerTestdata.FileWriterMock{}
	j := NewJob("file", "dir", "gradlew", "path", testdata.CmdFactoryMock{Name: "echo", Err: permissionErr}, fileWriterMock)

	go jobTestdata.WaitStatus(j)
	j.Run()

	assert.Len(t, j.Errors().GetAll(), 1)
}

func TestRunPermissionOutputErr(t *testing.T) {
	permissionErr := errors.New("asdhjaskdhqwe gradlew\": permission denied")
	otherErr := errors.New("WriteError")
	fileWriterMock := &writerTestdata.FileWriterMock{WriteErr: otherErr}

	j := NewJob("file", "dir", "gradlew", "path", testdata.CmdFactoryMock{Name: "bad-name", Err: permissionErr}, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.Len(t, j.Errors().GetAll(), 2)
}

func TestRun(t *testing.T) {
	fileContents := []byte("MakeDependenciesCmd\n")
	fileWriterMock := &writerTestdata.FileWriterMock{Contents: fileContents}
	cmdFactoryMock := testdata.CmdFactoryMock{Name: "echo"}
	j := NewJob("file", "dir", "gradlew", "path", cmdFactoryMock, fileWriterMock)

	go jobTestdata.WaitStatus(j)

	j.Run()

	assert.False(t, j.Errors().HasError())
	assert.Equal(t, fileContents, fileWriterMock.Contents)
}
