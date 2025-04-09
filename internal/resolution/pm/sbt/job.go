package sbt

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/debricked/cli/internal/resolution/job"
	"github.com/debricked/cli/internal/resolution/pm/maven"
	"github.com/debricked/cli/internal/resolution/pm/util"
)

const (
	executableNotFoundErrRegex = `executable file not found`
	pomGenerationErrRegex      = `Error occurred while processing command: makePom`
	sbtFileNotFoundErrRegex    = `not found: .*build\.sbt`
	nonParseableBuildErrRegex  = `Illegal character in build file`
	networkUnreachableErrRegex = `Connection timed out`
)

type Job struct {
	job.BaseJob
	cmdFactory      ICmdFactory
	buildService    IBuildService
	mavenPomService maven.IPomService
	mavenCmdFactory maven.ICmdFactory
}

func NewJob(file string, cmdFactory ICmdFactory, buildService IBuildService, mavenPomService maven.IPomService, mavenCmdFactory maven.ICmdFactory) *Job {
	return &Job{
		BaseJob:         job.NewBaseJob(file),
		cmdFactory:      cmdFactory,
		buildService:    buildService,
		mavenPomService: mavenPomService,
		mavenCmdFactory: mavenCmdFactory,
	}
}

func (j *Job) Run() {
	if err := j.parseBuildFile(); err != nil {
		return
	}

	if err := j.generatePomFile(); err != nil {
		return
	}

	pomFile, err := j.locatePomFile()
	if err != nil {
		return
	}

	pomXml, err := j.convertToPomXml(pomFile)
	if err != nil {
		return
	}

	if err := j.parseAndProcessWithMaven(pomXml); err != nil {
		return
	}
}

func (j *Job) parseBuildFile() error {
	status := "parsing SBT build file"
	j.SendStatus(status)

	file := j.GetFile()
	_, err := j.buildService.ParseBuildModules(file)
	if err != nil {
		doc := err.Error()
		if doc == "EOF" {
			doc = "This file doesn't contain valid SBT build content"
		}

		parsingError := util.NewPMJobError(err.Error())
		parsingError.SetStatus(status)
		parsingError.SetDocumentation(doc)
		j.Errors().Critical(parsingError)

		return err
	}

	return nil
}

func (j *Job) generatePomFile() error {
	file := j.GetFile()
	workingDirectory := filepath.Dir(filepath.Clean(file))
	cmd, err := j.cmdFactory.MakePomCmd(workingDirectory)
	if err != nil {
		j.handleError(util.NewPMJobError(err.Error()))

		return err
	}

	status := "generating Maven POM file"
	j.SendStatus(status)

	output, err := cmd.CombinedOutput()
	if err != nil {
		errContent := err.Error()
		if output != nil {
			errContent = string(output)
		}

		cmdErr := util.NewPMJobError(errContent)
		cmdErr.SetStatus(status)
		j.handleError(cmdErr)

		return err
	}

	return nil
}

func (j *Job) locatePomFile() (string, error) {
	file := j.GetFile()
	workingDirectory := filepath.Dir(filepath.Clean(file))
	status := "locating generated POM file"
	j.SendStatus(status)

	pomFile, err := j.buildService.FindPomFile(workingDirectory)
	if err != nil || pomFile == "" {
		errorMsg := "No pom file found in target directory"
		if err != nil {
			errorMsg = err.Error()
		}

		cmdErr := util.NewPMJobError(errorMsg)
		cmdErr.SetStatus(status)
		j.handleError(cmdErr)

		return "", err
	}

	return pomFile, nil
}

func (j *Job) convertToPomXml(pomFile string) (string, error) {
	file := j.GetFile()
	workingDirectory := filepath.Dir(filepath.Clean(file))
	status := "converting POM file to pom.xml"
	j.SendStatus(status)

	pomXml, err := j.buildService.RenamePomToXml(pomFile, workingDirectory)
	if err != nil {
		cmdErr := util.NewPMJobError(err.Error())
		cmdErr.SetStatus(status)
		j.handleError(cmdErr)

		return "", err
	}

	return pomXml, nil
}

func (j *Job) parseAndProcessWithMaven(pomXml string) error {
	status := "parsing Maven POM file"
	j.SendStatus(status)

	if err := j.parseMavenPom(pomXml); err != nil {
		return err
	}

	file := j.GetFile()
	workingDirectory := filepath.Dir(filepath.Clean(file))
	if err := j.createMavenDependencyGraph(workingDirectory, pomXml); err != nil {
		return err
	}

	status = fmt.Sprintf("processing dependencies with Maven resolver using %s", pomXml)
	j.SendStatus(status)

	return nil
}

func (j *Job) parseMavenPom(pomXml string) error {
	status := "parsing Maven POM file"
	j.SendStatus(status)

	_, err := j.mavenPomService.ParsePomModules(pomXml)
	if err != nil {
		doc := err.Error()
		if doc == "EOF" {
			doc = "This file doesn't contain valid XML"
		}

		parsingError := util.NewPMJobError(err.Error())
		parsingError.SetStatus(status)
		parsingError.SetDocumentation(doc)
		j.Errors().Critical(parsingError)

		return err
	}

	return nil
}

func (j *Job) createMavenDependencyGraph(workingDirectory string, pomXml string) error {
	status := "creating Maven dependency graph"
	j.SendStatus(status)

	cmd, err := j.mavenCmdFactory.MakeDependencyTreeCmd(workingDirectory)
	if err != nil {
		j.handleError(util.NewPMJobError(err.Error()))

		return err
	}

	output, err := cmd.Output()
	if err != nil {
		errContent := err.Error()
		if output != nil {
			errContent = string(output)
		}

		cmdErr := util.NewPMJobError(errContent)
		cmdErr.SetStatus(status)
		j.handleError(cmdErr)

		return err
	}

	return nil
}

func (j *Job) handleError(cmdError job.IError) {
	expressions := []string{
		executableNotFoundErrRegex,
		pomGenerationErrRegex,
		sbtFileNotFoundErrRegex,
		nonParseableBuildErrRegex,
		networkUnreachableErrRegex,
	}

	for _, expression := range expressions {
		regex := regexp.MustCompile(expression)
		if regex.MatchString(cmdError.Error()) {
			matches := regex.FindAllStringSubmatch(cmdError.Error(), -1)
			cmdError = j.addDocumentation(expression, matches, cmdError)
			j.Errors().Critical(cmdError)

			return
		}
	}

	j.Errors().Critical(cmdError)
}

func (j *Job) addDocumentation(expr string, matches [][]string, cmdError job.IError) job.IError {
	documentation := cmdError.Documentation()

	switch expr {
	case executableNotFoundErrRegex:
		documentation = j.GetExecutableNotFoundErrorDocumentation("SBT")
	case pomGenerationErrRegex:
		documentation = j.addPomGenerationErrorDocumentation(matches)
	case sbtFileNotFoundErrRegex:
		documentation = j.addSbtFileNotFoundErrorDocumentation(matches)
	case nonParseableBuildErrRegex:
		documentation = j.addNonParseableBuildErrorDocumentation(matches)
	case networkUnreachableErrRegex:
		documentation = j.addNetworkUnreachableErrorDocumentation()
	}

	cmdError.SetDocumentation(documentation)

	return cmdError
}

func (j *Job) addPomGenerationErrorDocumentation(matches [][]string) string {
	message := "Error occurred while generating the POM file"
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Failed to generate Maven POM file.",
			"SBT encountered an error during the makePom task.",
			"Error details:",
			message,
		}, " ")
}

func (j *Job) addSbtFileNotFoundErrorDocumentation(matches [][]string) string {
	message := "build.sbt file not found"
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"SBT configuration file not found.",
			"Please ensure that your project contains a valid build.sbt file.",
			"Error details:",
			message,
		}, " ")
}

func (j *Job) addNonParseableBuildErrorDocumentation(matches [][]string) string {
	message := "the build file for errors"
	if len(matches) > 0 && len(matches[0]) > 1 {
		message = matches[0][1]
	}

	return strings.Join(
		[]string{
			"Failed to parse SBT build file.",
			"Your build.sbt file contains syntax errors.",
			"Please check",
			message,
		}, " ")
}

func (j *Job) addNetworkUnreachableErrorDocumentation() string {
	return strings.Join(
		[]string{
			"We weren't able to retrieve one or more dependencies or plugins.",
			"Please check your Internet connection and try again.",
		}, " ")
}
