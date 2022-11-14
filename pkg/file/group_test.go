package file

import "testing"

func ExamplePrint() {
	format := Format{
		Regex:            "Regex",
		DocumentationUrl: "https://debricked.com/docs",
		LockFileRegexes:  []string{},
	}
	compiledFormat, _ := NewCompiledFormat(&format)
	g := NewGroup("package.json", compiledFormat, []string{"yarn.lock"})
	g.Print()
	// output:
	// package.json
	//  * yarn.lock
}

func TestHasFile(t *testing.T) {
	g := Group{
		FilePath:       "",
		CompiledFormat: nil,
		RelatedFiles:   []string{"yarn.lock"},
	}
	if g.HasFile() {
		t.Error("failed to assert that group had no dependency file")
	}
	g = Group{
		FilePath:       "package.json",
		CompiledFormat: nil,
		RelatedFiles:   []string{"yarn.lock"},
	}
	if !g.HasFile() {
		t.Error("failed to assert that group had a dependency file")
	}
}

func TestGetAllFiles(t *testing.T) {
	g := NewGroup("package.json", nil, []string{"yarn.lock"})
	if len(g.GetAllFiles()) != 2 {
		t.Error("failed to assert number of files")
	}
	g.RelatedFiles = []string{}
	if len(g.GetAllFiles()) != 1 {
		t.Error("failed to assert number of files")
	}
	g.FilePath = ""
	if len(g.GetAllFiles()) != 0 {
		t.Error("failed to assert number of files")
	}
}
