package language

import java "github.com/debricked/cli/internal/callgraph/language/java11"

type ILanguage interface {
	Name() string
	Version() string
}

func Languages() []ILanguage {
	return []ILanguage{
		java.NewLanguage(),
	}
}
