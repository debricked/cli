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
		sanitized = append(sanitized, sanitizeComposerPackage(entry))
	}

	return sanitized, nil
}

func sanitizeComposerPackage(pkg *util.OrderedObject) *util.OrderedObject {
	sanitized := util.NewOrderedObject()
	if name, isString := pkg.Values["name"].(string); isString {
		sanitized.Set("name", name)
	}
	if version, isString := pkg.Values["version"].(string); isString {
		sanitized.Set("version", sanitizeComposerVersion(version))
	}
	if requirements, isObject := pkg.Values["require"].(*util.OrderedObject); isObject {
		sanitizedRequirements := sanitizeComposerRequirements(requirements)
		if len(sanitizedRequirements.Keys) > 0 {
			sanitized.Set("require", sanitizedRequirements)
		}
	}

	return sanitized
}

func sanitizeComposerRequirements(requirements *util.OrderedObject) *util.OrderedObject {
	sanitized := util.NewOrderedObject()
	for _, name := range requirements.Keys {
		sanitized.Set(name, requirements.Values[name])
	}

	return sanitized
}

func sanitizeComposerVersion(version string) string {
	if strings.HasPrefix(version, "dev-") {
		return "dev-redacted"
	}

	return version
}
