package testdata

type ICallgraph interface {
	RunCallGraph() (string, error)
}

type CallgraphMock struct {
	RunCallGraphOutput string
	RunCallGraphError  error
}

func (cm CallgraphMock) RunCallGraph() (string, error) {
	return cm.RunCallGraphOutput, cm.RunCallGraphError

}
