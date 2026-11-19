package composer

import "github.com/debricked/cli/internal/sanitize/pm/util"

const Name = "composer"

type Pm struct {
	name string
}

func NewPm() Pm {
	return Pm{name: Name}
}

func (pm Pm) Name() string {
	return pm.name
}

func (Pm) LockFiles() map[string]util.FileSanitizer {
	return map[string]util.FileSanitizer{
		"composer.lock": sanitizeComposerLock,
	}
}

func (Pm) Manifests() map[string]util.FileSanitizer {
	return map[string]util.FileSanitizer{
		"composer.json": sanitizeComposerJSON,
	}
}
