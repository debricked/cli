package config

type IConfig interface {
	Language() string
	Path() string
	Arguments() []string
}

type Config struct {
	language  string
	execPath  string
	arguments []string
}

func NewConfig(language string, execPath string, arguments []string) Config {
	return Config{
		language,
		execPath,
		arguments,
	}
}

func (c Config) Language() string {
	return c.language
}

func (c Config) Path() string {
	return c.execPath
}

func (c Config) Arguments() []string {
	return c.arguments
}
