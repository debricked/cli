package util

import "slices"

// PackageJSONDependencyLists are the dependency lists of package.json. npm, yarn and pnpm lock
// files use the same names.
var PackageJSONDependencyLists = []string{
	"dependencies",
	"devDependencies",
	"optionalDependencies",
	"peerDependencies",
}

func IsPackageJSONDependencyList(key string) bool {
	return slices.Contains(PackageJSONDependencyLists, key)
}
