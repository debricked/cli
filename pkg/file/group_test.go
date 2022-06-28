package file

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
