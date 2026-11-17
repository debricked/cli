package composer

import (
	"strings"

	"github.com/debricked/cli/internal/sanitize/pm/util"
)

// sanitizeComposerLock keeps only the dependency data read by the backend parser.
func sanitizeComposerLock(content []byte, _ *util.FileContext) ([]byte, error) {
	lockFile, err := util.ParseJSONObject(content)
	if err != nil {
		return nil, util.ErrInvalidFormat
	}

	sanitized := util.NewOrderedObject()
	handled := false
	for _, key := range lockFile.Keys {
		if key != "packages" && key != "packages-dev" {
			continue
		}
		packages, isArray := lockFile.Values[key].([]any)
		if !isArray {
			return nil, util.ErrInvalidFormat
		}
		handled = true
		sanitizedPackages, sanitizeErr := sanitizeComposerPackages(packages)
		if sanitizeErr != nil {
			return nil, sanitizeErr
		}
		sanitized.Set(key, sanitizedPackages)
	}
	if !handled {
		return nil, util.ErrInvalidFormat
	}

	return util.EncodeJSON(sanitized)
}

func sanitizeComposerPackages(packages []any) ([]any, error) {
	sanitized := []any{}
	for _, pkg := range packages {
		entry, isObject := pkg.(*util.OrderedObject)
		if !isObject {
			return nil, util.ErrInvalidFormat
		}
		sanitizedEntry, err := sanitizeComposerPackage(entry)
		if err != nil {
			return nil, err
		}
		sanitized = append(sanitized, sanitizedEntry)
	}

	return sanitized, nil
}

// sanitizeComposerPackage keeps name and version, which make the component and its purl, and the
// names in require, which make the edges. Constraints are set to "*", as in composer.json.
// source and dist are dropped: they hold private URLs and paths.
func sanitizeComposerPackage(pkg *util.OrderedObject) (*util.OrderedObject, error) {
	name, isName := pkg.Values["name"].(string)
	version, isVersion := pkg.Values["version"].(string)
	if !isName || !isVersion {
		return nil, util.ErrInvalidFormat
	}
	sanitized := util.NewOrderedObject()
	sanitized.Set("name", name)
	sanitized.Set("version", sanitizeComposerVersion(version))

	value, present := pkg.Values["require"]
	if !present {
		return sanitized, nil
	}
	requirements, isObject := value.(*util.OrderedObject)
	if !isObject {
		return nil, util.ErrInvalidFormat
	}
	if len(requirements.Keys) > 0 {
		sanitized.Set("require", anyVersion(requirements))
	}

	return sanitized, nil
}

// sanitizeComposerVersion replaces branch versions. The parser reports every dev- version as
// empty, and branch names can hold ticket numbers and feature names.
func sanitizeComposerVersion(version string) string {
	if strings.HasPrefix(version, "dev-") {
		return "dev-redacted"
	}

	return version
}
