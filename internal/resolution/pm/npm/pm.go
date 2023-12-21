package npm

const Name = "npm"

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
		`package\.json$`,
	}
}
