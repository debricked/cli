// Package pm lists the package managers whose files sanitize can handle. Each one maps file names
// to the function that sanitizes that file type.
package pm

import (
	"github.com/debricked/cli/internal/sanitize/pm/npm"
	"github.com/debricked/cli/internal/sanitize/pm/nuget"
	"github.com/debricked/cli/internal/sanitize/pm/pnpm"
	"github.com/debricked/cli/internal/sanitize/pm/util"
	"github.com/debricked/cli/internal/sanitize/pm/yarn"
)

type IPm interface {
	Name() string
	// LockFiles are the dependency files the backend reads.
	LockFiles() map[string]util.FileSanitizer
	// Manifests are only written next to a lock file: the backend rejects them on their own.
	Manifests() map[string]util.FileSanitizer
}

func Pms() []IPm {
	return []IPm{
		npm.NewPm(),
		yarn.NewPm(),
		pnpm.NewPm(),
		nuget.NewPm(),
	}
}
