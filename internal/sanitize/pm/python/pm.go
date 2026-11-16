package python

import "github.com/debricked/cli/internal/sanitize/pm/util"

const Name = "python"

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
		"Pipfile.lock": sanitizePipfileLock,
		"poetry.lock":  sanitizePoetryLock,
		"uv.lock":      sanitizeUVLock,
	}
}

func (Pm) Manifests() map[string]util.FileSanitizer {
	return map[string]util.FileSanitizer{}
}
