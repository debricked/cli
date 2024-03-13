package golang

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path"
	"strings"

	"github.com/debricked/cli/internal/callgraph/cgexec"
	"github.com/debricked/cli/internal/callgraph/model"
	ioFs "github.com/debricked/cli/internal/io"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/callgraph/cha"
	"golang.org/x/tools/go/callgraph/rta"
	"golang.org/x/tools/go/callgraph/static"
	"golang.org/x/tools/go/callgraph/vta"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type ICallgraphBuilder interface {
	RunCallGraph() (string, error)
}

type CallgraphBuilder struct {
	filesystem       ioFs.IFileSystem
	workingDirectory string
	mainFile         string
	outputName       string
	ctx              cgexec.IContext
	cgModel          *model.CallGraph
	includeTests     bool
	algorithm        string
}

func NewCallgraphBuilder(
	workingDirectory string,
	mainFile string,
	outputName string,
	filesystem ioFs.IFileSystem,
	ctx cgexec.IContext,
) CallgraphBuilder {
	return CallgraphBuilder{
		workingDirectory: workingDirectory,
		mainFile:         mainFile,
		outputName:       outputName,
		filesystem:       filesystem,
		ctx:              ctx,
		cgModel:          model.NewCallGraph(),
		includeTests:     true,
		algorithm:        "cha",
	}
}

// Cleans up the symbol to remove pointer informaiton
// (*github.com/spf13/afero/mem.File).Open -> github.com/spf13/afero/mem.File.Open
func cleanSymbol(symbol string) string {
	symbol = strings.Replace(symbol, "(*", "", 1)
	symbol = strings.Replace(symbol, ").", ".", 1)

	return symbol
}

const (
	AlgorithmStatic = "static"
	AlgorithmCHA    = "cha"
	AlgorithmRTA    = "rta"
	AlgorithmVTA    = "vta"
)

// This is heavely inspired by: https://cs.opensource.google/go/x/tools/+/refs/tags/v0.19.0:cmd/callgraph/main.go
func (cg *CallgraphBuilder) constructCallGraph() error {
	cfg, args := cg.createPackageConfig()
	initial, err := packages.Load(cfg, args...)
	if err != nil {
		return fmt.Errorf("failed to load packages: %w", err)
	}

	prog, pkgs := cg.buildSSAFormProgram(initial)
	icg, err := cg.constructInternalCallGraph(prog, pkgs)
	if err != nil {
		return fmt.Errorf("failed to construct internal call graph: %w", err)
	}

	err = cg.outputCallGraph(icg, prog)
	if err != nil {
		return fmt.Errorf("failed to output call graph: %w", err)
	}

	err = setLineStartEnd(cg.cgModel)
	if err != nil {
		return fmt.Errorf("failed to set line start and end: %w", err)
	}

	return nil
}

func (cg *CallgraphBuilder) createPackageConfig() (*packages.Config, []string) {
	cfg := &packages.Config{
		Mode:  packages.NeedDeps | packages.NeedSyntax | packages.NeedTypesInfo | packages.NeedTypes | packages.NeedTypesSizes | packages.NeedImports | packages.NeedName | packages.NeedFiles | packages.NeedCompiledGoFiles,
		Tests: cg.includeTests,
		Dir:   cg.workingDirectory,
	}

	args := []string{cg.mainFile}

	return cfg, args
}

func (cg *CallgraphBuilder) buildSSAFormProgram(initial []*packages.Package) (*ssa.Program, []*ssa.Package) {
	mode := ssa.InstantiateGenerics // instantiate generics by default for soundness
	prog, pkgs := ssautil.AllPackages(initial, mode)
	prog.Build()

	return prog, pkgs
}

func (cg *CallgraphBuilder) constructInternalCallGraph(prog *ssa.Program, pkgs []*ssa.Package) (*callgraph.Graph, error) {
	var icg *callgraph.Graph
	var err error

	switch cg.algorithm {
	case AlgorithmStatic:
		icg = static.CallGraph(prog)

	case AlgorithmCHA:
		icg = cha.CallGraph(prog)

	case AlgorithmRTA:
		icg, err = cg.constructRTACallGraph(prog, pkgs)
		if err != nil {
			return nil, err
		}

	case AlgorithmVTA:
		icg = vta.CallGraph(ssautil.AllFunctions(prog), cha.CallGraph(prog))

	default:
		return nil, fmt.Errorf("unknown algorithm: %s", cg.algorithm)
	}

	icg.DeleteSyntheticNodes()

	return icg, nil
}

func (cg *CallgraphBuilder) constructRTACallGraph(prog *ssa.Program, pkgs []*ssa.Package) (*callgraph.Graph, error) {
	mains, err := mainPackages(pkgs)
	if err != nil {
		return nil, err
	}
	var roots []*ssa.Function
	for _, main := range mains {
		roots = append(roots, main.Func("init"), main.Func("main"))
	}
	rtares := rta.Analyze(roots, true)

	return rtares.CallGraph, nil
}

func isLibraryNode(filename string, pwd string) bool {
	if len(filename) > len(pwd) && filename[:len(pwd)] == pwd {
		return false
	}

	return true
}

func (cg *CallgraphBuilder) outputCallGraph(icg *callgraph.Graph, prog *ssa.Program) error {
	data := Edge{fset: prog.Fset}
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := callgraph.GraphVisitEdges(icg, func(edge *callgraph.Edge) error {
		data.position.Offset = -1
		data.callePosition.Offset = -1
		data.edge = edge
		data.Caller = edge.Caller.Func
		data.Callee = edge.Callee.Func
		callerSymbol := cleanSymbol(edge.Caller.Func.String())

		var callerNode *model.Node
		var ok bool
		if callerNode, ok = cg.cgModel.Nodes[callerSymbol]; !ok {
			callerFilename := data.CallerFilename()
			callerNode = cg.cgModel.AddNode(callerFilename, data.Caller.Name(), callerSymbol, isLibraryNode(callerFilename, pwd), -1, -1)
		}

		var calleeNode *model.Node
		calleeSymbol := cleanSymbol(edge.Callee.Func.String())
		if calleeNode, ok = cg.cgModel.Nodes[calleeSymbol]; !ok {
			calleeFilename := data.CalleeFilename()
			calleeNode = cg.cgModel.AddNode(calleeFilename, data.Callee.Name(), calleeSymbol, isLibraryNode(calleeFilename, pwd), -1, -1)
		}

		cg.cgModel.AddEdge(callerNode, calleeNode, data.CallLine())

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func setLineStartEnd(callgraph *model.CallGraph) error {
	fileNodeMap := map[string][]*model.Node{}

	for _, node := range callgraph.Nodes {
		if node.Filename != "" {
			fileNodeMap[node.Filename] = append(fileNodeMap[node.Filename], node)
		}
	}

	for filename, nodes := range fileNodeMap {
		if err := setLineStartEndFile(nodes, filename); err != nil {
			return err
		}
	}

	return nil
}

func setLineStartEndFile(nodes []*model.Node, filename string) error {
	if len(nodes) == 0 {
		return nil
	}

	fset := token.NewFileSet()
	fileNode, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	nameMap := map[string]*model.Node{}
	for _, node := range nodes {
		nameMap[node.Name] = node
	}

	ast.Inspect(fileNode, func(n ast.Node) bool {
		if fn, ok := n.(*ast.FuncDecl); ok {
			start := fset.Position(fn.Pos()).Line
			end := fset.Position(fn.End()).Line

			if node, ok := nameMap[fn.Name.Name]; ok {
				node.LineStart = start
				node.LineEnd = end
			}
		}

		return true
	})

	return nil
}

func (cg *CallgraphBuilder) RunCallGraph() (string, error) {

	err := cg.constructCallGraph()
	if err != nil {
		return "", err
	}

	cgOutputBytes, err := cg.cgModel.ToBytes()
	if err != nil {
		return "", err
	}

	outputFullPath := path.Join(cg.workingDirectory, cg.outputName)
	err = cg.filesystem.FsWriteFile(outputFullPath, cgOutputBytes, 0600)
	if err != nil {
		return "", err
	}

	return outputFullPath, nil
}

// mainPackages returns the main packages to analyze.
// Each resulting package is named "main" and has a main function.
func mainPackages(pkgs []*ssa.Package) ([]*ssa.Package, error) {
	var mains []*ssa.Package
	for _, p := range pkgs {
		if p != nil && p.Pkg.Name() == "main" && p.Func("main") != nil {
			mains = append(mains, p)
		}
	}
	if len(mains) == 0 {
		return nil, fmt.Errorf("no main packages")
	}

	return mains, nil
}

type Edge struct {
	Caller *ssa.Function
	Callee *ssa.Function

	edge          *callgraph.Edge
	fset          *token.FileSet
	position      token.Position // initialized lazily
	callePosition token.Position
}

func (e *Edge) pos() *token.Position {
	if e.position.Offset == -1 {
		e.position = e.fset.Position(e.edge.Pos()) // called lazily
	}

	return &e.position
}

func (e *Edge) calleePos() *token.Position {
	if e.callePosition.Offset == -1 {
		e.callePosition = e.fset.Position(e.Callee.Pos()) // called lazily
	}

	return &e.position
}

func (e *Edge) CallerFilename() string { return e.pos().Filename }
func (e *Edge) CalleeFilename() string { return e.calleePos().Filename }
func (e *Edge) CallLine() int          { return e.pos().Line }
