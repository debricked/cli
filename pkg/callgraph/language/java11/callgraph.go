package java

import (
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

//go:embed embeded/SootWrapper.jar
var jarCallGraph embed.FS

type Callgraph struct {
	cmdFactory       ICmdFactory
	workingDirectory string
	targetClasses    string
	targetDir        string
}

func (cg *Callgraph) runCallGraphWithSetup() error {
	jarFile, err := jarCallGraph.Open("embeded/SootWrapper.jar")
	if err != nil {
		return err
	}
	defer jarFile.Close()

	tempDir, err := ioutil.TempDir("", "jar")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	tempJarFile := filepath.Join(tempDir, "SootWrapper.jar")

	jarBytes, err := ioutil.ReadAll(jarFile)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(tempJarFile, jarBytes, 0644)
	if err != nil {
		return err
	}

	return cg.runCallGraph(tempJarFile)
}

func (cg *Callgraph) runCallGraph(callgraphJarPath string) error {
	cmd, err := cg.cmdFactory.MakeCallGraphGenerationCmd(callgraphJarPath, cg.workingDirectory, cg.targetClasses, cg.targetDir)
	fmt.Println(cmd)
	if err != nil {
		return err
	}
	fmt.Println("run command", cmd.Args)
	_, err = cmd.Output()

	fmt.Println("done running command", cmd.Args)
	return err
}
