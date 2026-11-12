package cocoapods

import "github.com/debricked/cli/internal/sanitize/pm/util"

const Name = "cocoapods"

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
		"Podfile.lock": sanitizePodfileLock,
	}
}

func (Pm) Manifests() map[string]util.FileSanitizer {
	return map[string]util.FileSanitizer{}
}
