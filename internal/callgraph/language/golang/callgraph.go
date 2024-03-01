package golang

import (
	"bytes"
	"path"
	"strconv"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	"github.com/debricked/cli/internal/callgraph/model"
	ioFs "github.com/debricked/cli/internal/io"
)

type ICallgraph interface {
	RunCallGraph() (string, error)
}

type Callgraph struct {
	cmdFactory       ICmdFactory
	filesystem       ioFs.IFileSystem
	workingDirectory string
	mainFile         string
	outputName       string
	ctx              cgexec.IContext
	cgModel          *model.CallGraph
}

func NewCallgraph(
	cmdFactory ICmdFactory,
	workingDirectory string,
	mainFile string,
	outputName string,
	filesystem ioFs.IFileSystem,
	ctx cgexec.IContext,
) Callgraph {
	return Callgraph{
		cmdFactory:       cmdFactory,
		workingDirectory: workingDirectory,
		mainFile:         mainFile,
		outputName:       outputName,
		filesystem:       filesystem,
		ctx:              ctx,
		cgModel:          model.NewCallGraph(),
	}
}

type IntermediateEdge struct {
	edges map[string][]string
}

func (ie *IntermediateEdge) AddEdge(parent, child string) {

	if _, ok := ie.edges[child]; !ok {
		ie.edges[child] = []string{}
	}

	ie.edges[child] = append(ie.edges[child], parent)

}

func (ie *IntermediateEdge) GetParents(child string) []string {
	return ie.edges[child]
}

func (cg *Callgraph) constructCallGraph(cgInput *[]byte) {

	uniqueNodeStrings := make(map[string]bool)
	intermediateEdge := &IntermediateEdge{
		edges: make(map[string][]string),
	}

	for _, line := range bytes.Split(*cgInput, []byte("\n")) {
		if len(line) == 0 {
			continue
		}
		// Clean out last and first character
		line = line[1 : len(line)-1]

		parts := bytes.Split(line, []byte("--->"))
		nodeParts := bytes.Split(parts[0], []byte(" "))
		filename := string(nodeParts[1])
		filename = filename[5:]

		symbol := string(nodeParts[0])
		symbolSplit := bytes.Split(nodeParts[0], []byte("."))
		name := string(symbolSplit[len(symbolSplit)-1])
		lineStart, _ := strconv.Atoi(string(nodeParts[2]))
		lineEnd, _ := strconv.Atoi(string(nodeParts[3]))
		node := cg.cgModel.AddNode(filename, name, symbol, false, lineStart, lineEnd)
		uniqueNodeStrings[node.Symbol] = true
		intermediateEdge.AddEdge(symbol, string(parts[1]))
	}

	for symbol := range uniqueNodeStrings {
		parents := intermediateEdge.GetParents(symbol)
		node := cg.cgModel.Nodes[symbol]

		for _, parent := range parents {
			parentNode := cg.cgModel.Nodes[parent]
			cg.cgModel.AddEdge(parentNode, node, 0)
		}
	}

}

func (cg *Callgraph) RunCallGraph() (string, error) {
	osCmd, err := cg.cmdFactory.MakeCallGraphGenerationCmd(cg.mainFile, cg.workingDirectory, cg.ctx)
	if err != nil {
		return "", err
	}

	cmd := cgexec.NewCommand(osCmd)
	err = cgexec.RunCommand(*cmd, cg.ctx)
	if err != nil {
		return "", err
	}

	output := cmd.GetStdOut()

	cgBytes := output.Bytes()
	cg.constructCallGraph(&cgBytes)
	cgOutputBytes, err := cg.cgModel.ToBytes()

	if err != nil {
		return "", err
	}
	outputFullPath := path.Join(cg.workingDirectory, cg.outputName)
	err = cg.filesystem.FsWriteFile(outputFullPath, cgOutputBytes, 0600)
	if err != nil {

		return "", err
	}
	_ = cmd.Wait()

	return outputFullPath, err
}
