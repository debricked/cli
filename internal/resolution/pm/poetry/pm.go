package poetry

const Name = "poetry"

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

func (_ Pm) Manifests() []string {
	return []string{
		"pyproject.toml",
	}
}
