package bundler

import "github.com/debricked/cli/internal/sanitize/pm/util"

const Name = "bundler"

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
		// Bundler also names lock files after their Gemfile, such as gems.Gemfile.lock.
		"*Gemfile.lock": sanitizeGemfileLock,
	}
}

func (Pm) Manifests() map[string]util.FileSanitizer {
	return map[string]util.FileSanitizer{}
}
