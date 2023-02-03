package gomod

const Name = "go"

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
		"go.mod",
	}
}
