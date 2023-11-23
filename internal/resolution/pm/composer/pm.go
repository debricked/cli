package composer

const Name = "composer"

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
		`composer\.json$`,
	}
}
