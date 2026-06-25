package language

import "github.com/debricked/cli/internal/callgraph/language/java"

type ILanguage interface {
	Name() string
	Version() string
}

func Languages() []ILanguage {
	return []ILanguage{
		java.NewLanguage(),
	}
}
