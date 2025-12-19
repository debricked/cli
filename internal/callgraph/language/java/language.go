package java

const Name = "java"
const StandardVersion = "11"

type Language struct {
	name    string
	version string
}

func NewLanguage() Language {
	return Language{
		name:    Name,
		version: StandardVersion,
	}
}

func (language Language) Name() string {
	return language.name
}

func (language Language) Version() string {
	return language.version
}
