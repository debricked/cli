package cargo

import "github.com/debricked/cli/internal/sanitize/pm/util"

const Name = "cargo"

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

func (Pm) LockFiles() map[string]util.FileSanitizer {
	return map[string]util.FileSanitizer{
		"Cargo.lock": sanitizeCargoLock,
	}
}

func (Pm) Manifests() map[string]util.FileSanitizer {
	return map[string]util.FileSanitizer{}
}
