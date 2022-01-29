package file

func ExamplePrint() {
	f := Format{
		Regex:            "Regex",
		DocumentationUrl: "https://debricked.com/docs",
		LockFileRegexes:  []string{},
	}
	g := NewGroup("package.json", &f, []string{"yarn.lock"})
	g.Print()
	// output:
	// package.json
	//  * yarn.lock
}
