package testdata

type ICallgraph interface {
	RunCallGraphWithSetup() error
	RunCallGraph(callgraphJarPath string) error
}

type CallgraphMock struct {
	RunCallGraphWithSetupError error
	RunCallGraphError          error
}

func (cm CallgraphMock) RunCallGraphWithSetup() error {
	return cm.RunCallGraphWithSetupError
}

func (cm CallgraphMock) RunCallGraph(callgraphJarPath string) error {
	return cm.RunCallGraphError

}
