package golang

const Name = "golang"
const StandardVersion = "1"

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
