package testdata

type MockGreeter struct{}

func (MockGreeter) Greeting(string) string {
	return "hello"
}
