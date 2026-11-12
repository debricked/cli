package bazel

import "github.com/debricked/cli/internal/sanitize/pm/util"

const Name = "bazel"

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
		// rules_jvm_external names the file after the repository, such as maven_install.json.
		"*_install.json": sanitizeInstallJSON,
	}
}

func (Pm) Manifests() map[string]util.FileSanitizer {
	return map[string]util.FileSanitizer{}
}
