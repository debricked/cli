package java

import (
	"regexp"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	ioFs "github.com/debricked/cli/internal/io"
)

type ICallgraph interface {
	RunCallGraphWithSetup() error
	RunCallGraph(callgraphJarPath string) error
}

type Callgraph struct {
	cmdFactory       ICmdFactory
	filesystem       ioFs.IFileSystem
	archive          ioFs.IArchive
	workingDirectory string
	targetClasses    []string
	targetDir        string
	outputName       string
	ctx              cgexec.IContext
	sootHandler      ISootHandler
}

func NewCallgraph(
	cmdFactory ICmdFactory,
	workingDirectory string,
	targetClasses []string,
	targetDir string,
	outputName string,
	filesystem ioFs.IFileSystem,
	archive ioFs.IArchive,
	ctx cgexec.IContext,
	sootHandler ISootHandler,
) Callgraph {
	return Callgraph{
		cmdFactory:       cmdFactory,
		workingDirectory: workingDirectory,
		targetClasses:    targetClasses,
		targetDir:        targetDir,
		outputName:       outputName,
		filesystem:       filesystem,
		archive:          archive,
		ctx:              ctx,
		sootHandler:      sootHandler,
	}
}

func (cg *Callgraph) RunCallGraphWithSetup() error {
	version, err := cg.javaVersion(".")
	if err != nil {
		return err
	}

	jarFile, err := cg.sootHandler.GetSootWrapper(version, cg.filesystem, cg.archive)
	if err != nil {

		return err
	}

	err = cg.RunCallGraph(jarFile)

	return err
}

func (cg *Callgraph) javaVersion(path string) (string, error) {
	osCmd, err := cg.cmdFactory.MakeJavaVersionCmd(path, cg.ctx)
	if err != nil {

		return "", err
	}

	cmd := cgexec.NewCommand(osCmd)
	err = cgexec.RunCommand(*cmd, cg.ctx)
	javaVersionRegex := regexp.MustCompile(`\b(\d+)\.\d+\.\d+\b`)
	match := javaVersionRegex.FindStringSubmatch(cmd.GetStdOut().String())
	if len(match) > 1 {
		return match[1], nil
	}

	return "", err
}

func (cg *Callgraph) RunCallGraph(callgraphJarPath string) error {
	osCmd, err := cg.cmdFactory.MakeCallGraphGenerationCmd(
		callgraphJarPath,
		cg.workingDirectory,
		cg.targetClasses,
		cg.targetDir,
		cg.outputName,
		cg.ctx,
	)
	if err != nil {

		return err
	}

	cmd := cgexec.NewCommand(osCmd)
	err = cgexec.RunCommand(*cmd, cg.ctx)

	return err
}
