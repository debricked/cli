package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeToBytes(t *testing.T) {
	node := &Node{
		Filename:          "file.go",
		Name:              "Node",
		Symbol:            "symbol",
		IsApplicationNode: false,
		IsStdLibNode:      true,
		LineStart:         1,
		LineEnd:           10,
		Parents:           []Edge{},
	}

	bytes := node.ToBytes()
	assert.NotNil(t, bytes)
	assert.Equal(t, "[\"symbol\", false, true, \"Node\", \"file.go\", 1, 10, []]", string(bytes))

	node.IsApplicationNode = true
	bytes = node.ToBytes()
	assert.NotNil(t, bytes)
	assert.Equal(t, "[\"symbol\", true, true, \"Node\", \"file.go\", 1, 10, []]", string(bytes))

	node.Parents = append(node.Parents, Edge{Parent: node, CallLine: 10})
	bytes = node.ToBytes()
	assert.NotNil(t, bytes)
	assert.Equal(t, "[\"symbol\", true, true, \"Node\", \"file.go\", 1, 10, [[\"symbol\", 10, \"file.go\"]]]", string(bytes))

}

func TestNewCallGraph(t *testing.T) {
	cg := NewCallGraph()
	assert.NotNil(t, cg)
	assert.Equal(t, CURRENT_VERSION, cg.Version)
	assert.NotNil(t, cg.Nodes)
}

func TestAddNode(t *testing.T) {
	cg := NewCallGraph()
	node := cg.AddNode("file.go", "Node", "symbol", false, false, 1, 10)
	assert.NotNil(t, node)
	assert.Equal(t, "file.go", node.Filename)
	assert.Equal(t, "Node", node.Name)
	assert.Equal(t, "symbol", node.Symbol)
	assert.Equal(t, false, node.IsApplicationNode)
	assert.Equal(t, 1, node.LineStart)
	assert.Equal(t, 10, node.LineEnd)
	assert.NotNil(t, node.Parents)

	// Add the same node again
	node2 := cg.AddNode("file.go", "Node", "symbol", false, false, 1, 10)
	assert.Equal(t, node, node2)
}

func TestAddEdge(t *testing.T) {
	cg := NewCallGraph()
	parent := cg.AddNode("file.go", "Parent", "symbol1", false, false, 1, 10)
	child := cg.AddNode("file.go", "Child", "symbol2", false, false, 11, 20)
	cg.AddEdge(parent, child, 10)
	assert.Equal(t, 1, len(child.Parents))
	assert.Equal(t, parent, child.Parents[0].Parent)
	assert.Equal(t, 10, child.Parents[0].CallLine)
}

func TestGetNode(t *testing.T) {
	cg := NewCallGraph()
	node := cg.AddNode("file.go", "Node", "symbol", false, false, 1, 10)
	gotNode := cg.GetNode("symbol")
	assert.Equal(t, node, gotNode)
}

func TestNodeCount(t *testing.T) {
	cg := NewCallGraph()
	cg.AddNode("file.go", "Node1", "symbol1", false, false, 1, 10)
	cg.AddNode("file.go", "Node2", "symbol2", false, false, 11, 20)
	assert.Equal(t, 2, cg.NodeCount())
}

func TestEdgeCount(t *testing.T) {
	cg := NewCallGraph()
	parent := cg.AddNode("file.go", "Parent", "symbol1", false, false, 1, 10)
	child1 := cg.AddNode("file.go", "Child1", "symbol2", false, false, 11, 20)
	child2 := cg.AddNode("file.go", "Child2", "symbol3", false, false, 21, 30)
	cg.AddEdge(parent, child1, 10)
	cg.AddEdge(parent, child2, 20)
	assert.Equal(t, 2, cg.EdgeCount())
}

func TestCallgRaphToBytes(t *testing.T) {
	cg := NewCallGraph()
	cg.AddNode("file.go", "Node1", "symbol1", false, false, 1, 10)
	cg.AddNode("file.go", "Node2", "symbol2", false, false, 11, 20)
	cg.AddEdge(cg.GetNode("symbol1"), cg.GetNode("symbol2"), 10)
	cg.AddEdge(cg.GetNode("symbol1"), cg.GetNode("symbol2"), 20)

	bytes, err := cg.ToBytes()
	assert.Nil(t, err)
	assert.NotNil(t, bytes)
	assert.Equal(t, "{\"version\": \"5\", \"data\": [[\"symbol1\", false, false, \"Node1\", \"file.go\", 1, 10, []],[\"symbol2\", false, false, \"Node2\", \"file.go\", 11, 20, [[\"symbol1\", 10, \"file.go\"], [\"symbol1\", 20, \"file.go\"]]]]}", string(bytes))
}
