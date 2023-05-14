package java

import (
	"embed"
	"fmt"
	"os"
	"os/exec"
	"path"

	conf "github.com/debricked/cli/pkg/callgraph/config"
	"github.com/debricked/cli/pkg/callgraph/job"
	"github.com/debricked/cli/pkg/io/finder"
	ioWriter "github.com/debricked/cli/pkg/io/writer"
)

const (
	maven  = "maven"
	gradle = "gradle"
)

type Job struct {
	job.BaseJob
	cmdFactory ICmdFactory
	config     conf.IConfig
}

func NewJob(dir string, files []string, cmdFactory ICmdFactory, writer ioWriter.IFileWriter, config conf.IConfig) *Job {
	return &Job{
		BaseJob:    job.NewBaseJob(dir, files),
		cmdFactory: cmdFactory,
		config:     config,
	}
}

//go:embed embeded/gradle-script.groovy
var gradleInitScript embed.FS

func (j *Job) Run() {
	fmt.Println("ENTERED RUN")
	workingDirectory := j.GetDir()
	fmt.Println("Files:", j.GetFiles())
	targetClasses := j.GetFiles()[0]
	dependencyDir := ".debrickedTmpFolder"
	targetDir := path.Join(workingDirectory, dependencyDir)
	pmConfig := j.config.Kwargs()["pm"]

	// If folder doesn't exist, copy dependencies
	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		var cmd *exec.Cmd
		if pmConfig == gradle {
			targetGradlew := path.Join(workingDirectory, "gradlew")
			gradlew := "gradle"
			if _, err := os.Stat(targetGradlew); os.IsExist(err) {
				gradlew = targetGradlew
			}

			groovyFilePath := path.Join(workingDirectory, ".debrickedGroovyScript.groovy")
			ish := finder.NewScriptHandler(groovyFilePath, "embeded/gradle-script.groovy", ioWriter.FileWriter{})
			ish.WriteInitFile()

			cmd, err = j.cmdFactory.MakeGradleCopyDependenciesCmd(workingDirectory, gradlew, groovyFilePath)
		} else {
			cmd, err = j.cmdFactory.MakeMvnCopyDependenciesCmd(workingDirectory, targetDir)
		}
		fmt.Println("Copying relevant jars to target folder", targetDir, cmd.Args)
		if err != nil {
			j.Errors().Critical(err)

			return
		}
		_, err = cmd.Output()

		if err != nil {
			j.Errors().Critical(err)

			return
		}
	}

	fmt.Println("Generating callgraph!")
	j.SendStatus("generating call graph")
	callgraph := Callgraph{
		cmdFactory:       j.cmdFactory,
		workingDirectory: workingDirectory,
		targetClasses:    targetClasses,
		targetDir:        targetDir,
	}
	err := callgraph.runCallGraphWithSetup()

	if err != nil {
		j.Errors().Critical(err)
	}
}
