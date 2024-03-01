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
	Symbol   string
	CallLine int
}

type IntermediateEdges struct {
	edges map[string][]IntermediateEdge
}

func (ie *IntermediateEdges) AddEdge(parent, child string, line int) {

	if _, ok := ie.edges[child]; !ok {
		ie.edges[child] = make([]IntermediateEdge, 0)
	}

	ie.edges[child] = append(ie.edges[child], IntermediateEdge{Symbol: parent, CallLine: line})

}

func (ie *IntermediateEdges) GetParents(child string) []IntermediateEdge {
	return ie.edges[child]
}

func (cg *Callgraph) constructCallGraph(cgInput *[]byte) {

	uniqueNodeStrings := make(map[string]bool)
	intermediateEdge := &IntermediateEdges{
		edges: make(map[string][]IntermediateEdge),
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
		callLine, _ := strconv.Atoi(string(nodeParts[2]))
		lineStart := -1 //, _ := strconv.Atoi(string(nodeParts[2]))
		lineEnd := -1
		node := cg.cgModel.AddNode(filename, name, symbol, false, lineStart, lineEnd)
		uniqueNodeStrings[node.Symbol] = true
		intermediateEdge.AddEdge(symbol, string(parts[1]), callLine)
	}

	for symbol := range uniqueNodeStrings {
		parents := intermediateEdge.GetParents(symbol)
		node := cg.cgModel.Nodes[symbol]

		for _, parent := range parents {
			parentNode := cg.cgModel.Nodes[parent.Symbol]
			cg.cgModel.AddEdge(parentNode, node, parent.CallLine)
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
