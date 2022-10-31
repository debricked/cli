package file

import "path/filepath"

func DefaultExclusions() []string {
	return []string{
		filepath.Join("**", "node_modules", "**"),
		filepath.Join("**", "vendor", "**"),
		filepath.Join("**", ".git", "**"),
	}
}
