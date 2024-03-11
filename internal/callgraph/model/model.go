package model

import (
	"bytes"
	"fmt"
)

const CURRENT_VERSION = "5"

type Edge struct {
	Parent   *Node
	CallLine int
}

type Node struct {
	Filename      string
	Name          string
	Symbol        string
	IsLibraryNode bool
	LineStart     int
	LineEnd       int
	Parents       []Edge
}

func (n *Node) ToBytes() []byte {
	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("['%s', ", n.Symbol))

	if n.IsLibraryNode {
		buffer.WriteString("True, True, ")
	} else {
		buffer.WriteString("False, False, ")
	}

	buffer.WriteString(fmt.Sprintf("'%s', '%s', %d, %d, [", n.Name, n.Filename, n.LineStart, n.LineEnd))

	for _, parent := range n.Parents {
		buffer.WriteString(fmt.Sprintf("['%s', %d, '%s'], ", parent.Parent.Symbol, parent.CallLine, parent.Parent.Filename))
	}

	buffer.WriteString("]]")

	return buffer.Bytes()
}

type CallGraph struct {
	Nodes   map[string]*Node
	Version string
}

func NewCallGraph() *CallGraph {
	return &CallGraph{
		Nodes:   map[string]*Node{},
		Version: CURRENT_VERSION,
	}
}

func (cg *CallGraph) AddNode(filename, name, symbol string, isLibraryNode bool, lineStart, lineEnd int) *Node {

	if cg.Nodes == nil {
		cg.Nodes = make(map[string]*Node)
	}

	if node, ok := cg.Nodes[symbol]; ok {
		return node
	}

	node := &Node{
		Filename:      filename,
		Name:          name,
		Symbol:        symbol,
		IsLibraryNode: isLibraryNode,
		LineStart:     lineStart,
		LineEnd:       lineEnd,
		Parents:       []Edge{},
	}

	cg.Nodes[symbol] = node

	return node
}

func (cg *CallGraph) AddEdge(parent, child *Node, callLine int) {
	callee := Edge{
		Parent:   parent,
		CallLine: callLine,
	}
	child.Parents = append(child.Parents, callee)
}

func (cg *CallGraph) GetNode(symbol string) *Node {
	return cg.Nodes[symbol]
}

func (cg *CallGraph) ToBytes() ([]byte, error) {
	output := []byte{}

	output = append(output, []byte("[")...)

	for _, node := range cg.Nodes {
		output = append(output, node.ToBytes()...)
		output = append(output, []byte(",")...)
	}

	if len(output) > 0 {
		output = output[:len(output)-1]
	}

	output = append(output, []byte("]")...)

	output = append([]byte("{\"version\": \""+cg.Version+"\", \"data\": "), output...)
	output = append(output, []byte("}")...)

	return output, nil
}

func (cg *CallGraph) NodeCount() int {
	return len(cg.Nodes)
}

func (cg *CallGraph) EdgeCount() int {
	count := 0
	for _, node := range cg.Nodes {
		count += len(node.Parents)
	}
	return count
}
