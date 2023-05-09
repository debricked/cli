package config

type IConfig interface {
	Language() string
	Args() []string
	Kwargs() map[string]string
}

type Config struct {
	language string
	args     []string
	kwargs   map[string]string
}

func NewConfig(language string, args []string, kwargs map[string]string) Config {
	return Config{
		language,
		args,
		kwargs,
	}
}

func (c Config) Language() string {
	return c.language
}

func (c Config) Args() []string {
	return c.args
}

func (c Config) Kwargs() map[string]string {
	return c.kwargs
}
