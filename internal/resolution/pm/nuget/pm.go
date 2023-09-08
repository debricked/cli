package nuget

const Name = "nuget"
const CsprojRegex = `\.csproj$`
const PackagesConfigRegex = `packages\.config$`

type Pm struct {
	name string
}

func NewPm() Pm {
	return Pm{
		name: Name,
	}
}

func (pm Pm) Name() string {
	return pm.name
}

func (Pm) Manifests() []string {
	return []string{
		CsprojRegex,
		PackagesConfigRegex,
	}
}
