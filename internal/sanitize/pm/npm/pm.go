package npm

import "github.com/debricked/cli/internal/sanitize/pm/util"

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

func (Pm) LockFiles() map[string]util.FileSanitizer {
	return map[string]util.FileSanitizer{
		"package-lock.json": sanitizePackageLock,
	}
}

func (Pm) Manifests() map[string]util.FileSanitizer {
	return map[string]util.FileSanitizer{
		// package.json is listed once, here, although yarn and pnpm projects use it too.
		"package.json": sanitizePackageJSON,
	}
}
