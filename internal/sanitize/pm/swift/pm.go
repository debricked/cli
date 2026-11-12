package swift

import "github.com/debricked/cli/internal/sanitize/pm/util"

const Name = "swift"

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
		"Package.resolved": sanitizePackageResolved,
	}
}

func (Pm) Manifests() map[string]util.FileSanitizer {
	return map[string]util.FileSanitizer{
		// Package.swift is source code, so it is never copied. The backend only uses it to set
		// the scope to production, which doesn't change the result.
	}
}
