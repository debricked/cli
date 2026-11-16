// Package testdata has helpers for the tests of the sanitize package managers.
package testdata

import "github.com/debricked/cli/internal/sanitize/pm/util"

// NewContext returns a context for the file at path, in the first project of a fresh run, and the
// placeholders of that run.
func NewContext(path string) (*util.FileContext, *util.Placeholders) {
	all := &util.Placeholders{}

	return util.NewFileContext(util.NewProjectScope(all, 1), path), all
}
