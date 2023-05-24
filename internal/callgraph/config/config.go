package config

type IConfig interface {
	Language() string
	Args() []string
	Kwargs() map[string]string
	Build() bool
	PackageManager() string
}

type Config struct {
	language       string
	args           []string
	kwargs         map[string]string
	build          bool
	packageManager string
}

func NewConfig(language string, args []string, kwargs map[string]string, build bool, packageManager string) Config {
	return Config{
		language,
		args,
		kwargs,
		build,
		packageManager,
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

func (c Config) Build() bool {
	return c.build
}

func (c Config) PackageManager() string {
	return c.packageManager
}
