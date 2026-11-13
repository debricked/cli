package pnpm

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/debricked/cli/internal/sanitize/pm/util"
)

const (
	pnpmImporterFormat  = "redacted-%d"
	pnpmPackageFormat   = "redacted-%d"
	pnpmLocalDotFormat  = "file:./redacted-%d"
	pnpmLocalFormat     = "file:redacted-%d"
	pnpmRootImporterKey = "."
)

var (
	// pnpmSlashKeyRegex matches registry keys up to v6: /name/1.0.0, /@scope/name@1.0.0.
	pnpmSlashKeyRegex = regexp.MustCompile(`^/(?:@[^/@]+/)?[^/@]+[/@][\w.+~@-]+$`)
	// pnpmNamedKeyRegex matches keys from v9, where the version or source follows the name.
	pnpmNamedKeyRegex = regexp.MustCompile(`^((?:@[^/@:]+/)?[^/@:]+)@(.+)$`)
	// pnpmVersionRegex matches versions, including the _peer@1.0.0 suffix of older lock files.
	pnpmVersionRegex = regexp.MustCompile(`^\d[\w.+~@-]*$`)
)

// pnpmLockSanitizer builds its output from new values only and never copies nodes of the source
// document. That way an unknown key or a comment can't reach the output.
type pnpmLockSanitizer struct {
	context  *util.FileContext
	baseKeys map[string]string
	usedKeys map[string]int
}

func sanitizePnpmLock(content []byte, context *util.FileContext) ([]byte, error) {
	root, err := util.ParseYAMLMapping(content)
	if err != nil {
		return nil, err
	}

	sanitizer := &pnpmLockSanitizer{context: context, baseKeys: map[string]string{}, usedKeys: map[string]int{}}

	// Package keys get their placeholders first, so references to them can be replaced wherever
	// they appear.
	if packages := util.YAMLValue(root, "packages"); packages != nil && packages.Kind == yaml.MappingNode {
		if err = sanitizer.registerKeys(packages); err != nil {
			return nil, util.ErrInvalidFormat
		}
	}

	sanitized, err := sanitizer.sanitizeRoot(root)
	if err != nil {
		return nil, util.ErrInvalidFormat
	}

	return util.EncodeYAML(sanitized), nil
}

// sanitizeRoot fails if no known section is left, because then the file isn't a pnpm lock file.
func (s *pnpmLockSanitizer) sanitizeRoot(root *yaml.Node) (*util.OrderedObject, error) {
	sanitized := util.NewOrderedObject()
	var err error
	util.ForEachYAMLPair(root, func(key string, value *yaml.Node) {
		if err != nil {
			return
		}
		var section any
		var found bool
		if section, found, err = s.sanitizeSection(key, value); found {
			sanitized.Set(key, section)
		}
	})
	if err == nil && len(sanitized.Keys) == 0 {
		err = util.ErrInvalidFormat
	}

	return sanitized, err
}

// sanitizeSection returns found = false for a key that isn't allow-listed, or a section that is
// empty in the source file. Both are left out of the output.
func (s *pnpmLockSanitizer) sanitizeSection(key string, value *yaml.Node) (any, bool, error) {
	switch key {
	case "lockfileVersion":
		if value.Kind != yaml.ScalarNode {
			return nil, false, util.ErrInvalidFormat
		}

		return pnpmLockfileVersion(value), true, nil
	case "importers":
		return optionalSection(s.sanitizeImporters(value))
	case "packages":
		return optionalSection(s.sanitizeEntries(value, pnpmPackageAllowList))
	case "snapshots":
		return optionalSection(s.sanitizeEntries(value, pnpmSnapshotAllowList))
	case "dependencies", "devDependencies", "optionalDependencies":
		// A single-project lock file up to v6 has these in place of importers.
		return sanitizeNames(value), true, nil
	default:
		return nil, false, nil
	}
}

// errEmptySection marks a section that is null in the source file, such as `importers:` with
// nothing under it.
var errEmptySection = errors.New("section is empty")

func optionalSection(section *util.OrderedObject, err error) (any, bool, error) {
	switch {
	case errors.Is(err, errEmptySection):
		return nil, false, nil
	case err != nil:
		return nil, false, err
	default:
		return section, true, nil
	}
}

func pnpmLockfileVersion(value *yaml.Node) any {
	if value.Tag == "!!float" || value.Tag == "!!int" {
		return util.RawYAMLScalar(value.Value)
	}

	return value.Value
}

func (s *pnpmLockSanitizer) sanitizeImporters(importers *yaml.Node) (*util.OrderedObject, error) {
	if util.IsYAMLNull(importers) {
		return nil, errEmptySection
	}
	if importers.Kind != yaml.MappingNode {
		return nil, util.ErrInvalidFormat
	}

	sanitized := util.NewOrderedObject()
	count := 0
	util.ForEachYAMLPair(importers, func(path string, importer *yaml.Node) {
		key := path
		if path != pnpmRootImporterKey {
			count++
			key = fmt.Sprintf(pnpmImporterFormat, count)
			s.context.RecordPlaceholder(key, path)
		}
		sanitized.Set(key, sanitizeImporter(importer))
	})

	return sanitized, nil
}

// sanitizeImporter keeps names only: the values often hold internal URLs and aren't read.
func sanitizeImporter(importer *yaml.Node) *util.OrderedObject {
	sanitized := util.NewOrderedObject()
	util.ForEachYAMLPair(importer, func(listName string, list *yaml.Node) {
		if !util.IsPackageJSONDependencyList(listName) {
			return
		}
		sanitized.Set(listName, sanitizeNames(list))
	})

	return sanitized
}

func sanitizeNames(list *yaml.Node) *util.OrderedObject {
	names := util.NewOrderedObject()
	util.ForEachYAMLPair(list, func(name string, _ *yaml.Node) {
		names.Set(name, util.NewOrderedObject())
	})

	return names
}

type pnpmAllowList struct {
	scalars map[string]bool
	lists   map[string]bool
}

var (
	pnpmPackageAllowList = pnpmAllowList{
		scalars: map[string]bool{"name": true, "version": true},
		lists:   map[string]bool{"dependencies": true, "optionalDependencies": true, "peerDependencies": true},
	}
	pnpmSnapshotAllowList = pnpmAllowList{
		lists: map[string]bool{"dependencies": true, "optionalDependencies": true},
	}
)

func (s *pnpmLockSanitizer) sanitizeEntries(entries *yaml.Node, allowed pnpmAllowList) (*util.OrderedObject, error) {
	if util.IsYAMLNull(entries) {
		return nil, errEmptySection
	}
	if entries.Kind != yaml.MappingNode {
		return nil, util.ErrInvalidFormat
	}

	sanitized := util.NewOrderedObject()
	var err error
	util.ForEachYAMLPair(entries, func(key string, entry *yaml.Node) {
		if err != nil {
			return
		}
		var maskedKey string
		var sanitizedEntry *util.OrderedObject
		if maskedKey, err = s.maskKey(key); err != nil {
			return
		}
		if sanitizedEntry, err = s.sanitizeEntry(entry, allowed); err != nil {
			return
		}
		// Two entries ending up with one key would silently lose a package.
		if _, exists := sanitized.Values[maskedKey]; exists {
			err = util.ErrInvalidFormat

			return
		}
		sanitized.Set(maskedKey, sanitizedEntry)
	})

	return sanitized, err
}

func (s *pnpmLockSanitizer) sanitizeEntry(entry *yaml.Node, allowed pnpmAllowList) (*util.OrderedObject, error) {
	sanitized := util.NewOrderedObject()
	var err error
	util.ForEachYAMLPair(entry, func(key string, value *yaml.Node) {
		if err != nil {
			return
		}
		switch {
		case allowed.scalars[key] && value.Kind == yaml.ScalarNode && !util.IsYAMLNull(value):
			if key == "version" {
				sanitized.Set(key, util.MaskValue(value.Value))
			} else {
				sanitized.Set(key, value.Value)
			}
		case allowed.lists[key]:
			var list *util.OrderedObject
			if list, err = s.sanitizeDependencyList(value); err == nil {
				sanitized.Set(key, list)
			}
		}
	})

	return sanitized, err
}

func (s *pnpmLockSanitizer) sanitizeDependencyList(list *yaml.Node) (*util.OrderedObject, error) {
	sanitized := util.NewOrderedObject()
	var err error
	util.ForEachYAMLPair(list, func(name string, value *yaml.Node) {
		if err != nil || value.Kind != yaml.ScalarNode || util.IsYAMLNull(value) {
			return
		}
		var masked string
		if masked, err = s.maskReference(name, value.Value); err == nil {
			sanitized.Set(name, masked)
		}
	})

	return sanitized, err
}

func (s *pnpmLockSanitizer) registerKeys(packages *yaml.Node) error {
	var err error
	util.ForEachYAMLPair(packages, func(key string, _ *yaml.Node) {
		if err == nil {
			_, err = s.maskKey(key)
		}
	})

	return err
}

// maskKey masks a package key. The part in brackets lists peer dependencies, and each of those
// is a package key too.
func (s *pnpmLockSanitizer) maskKey(key string) (string, error) {
	base, peers := splitPnpmPeers(key)
	maskedPeers, err := s.maskPeers(peers)
	if err != nil {
		return "", err
	}

	return s.maskBaseKey(base) + maskedPeers, nil
}

func (s *pnpmLockSanitizer) maskBaseKey(base string) string {
	if masked, found := s.baseKeys[base]; found {
		return masked
	}
	masked := s.uniqueKey(s.computeBaseKey(base))
	s.baseKeys[base] = masked

	return masked
}

func (s *pnpmLockSanitizer) computeBaseKey(base string) string {
	if pnpmSlashKeyRegex.MatchString(base) {
		return base
	}
	if rest, isLocal := util.CutLocalProtocol(base); isLocal {
		format := pnpmLocalFormat
		if strings.HasPrefix(rest, ".") {
			format = pnpmLocalDotFormat
		}

		return s.context.PathPlaceholder(format, base)
	}
	// An unrecognised source means this isn't a name and a version after all. For example,
	// git@github.com+org/repo/sha would otherwise give a package called git.
	if match := pnpmNamedKeyRegex.FindStringSubmatch(base); match != nil {
		if masked := maskPnpmVersion(match[2]); masked != util.MaskedUnknown {
			return match[1] + "@" + masked
		}
	}

	// Not a form known to be safe, such as a host and a path: nothing of it is kept.
	return s.context.PathPlaceholder(pnpmPackageFormat, base)
}

// uniqueKey is needed because masking can give two sources of one package the same key, and a
// YAML mapping can't hold both.
func (s *pnpmLockSanitizer) uniqueKey(key string) string {
	s.usedKeys[key]++
	count := s.usedKeys[key]
	if count == 1 {
		return key
	}
	if strings.HasSuffix(key, util.MaskedUnknown) {
		return fmt.Sprintf("%s-%d", key, count)
	}

	return fmt.Sprintf("%s#%s-%d", key, util.MaskedUnknown, count)
}

func (s *pnpmLockSanitizer) maskPeers(peers string) (string, error) {
	var masked strings.Builder
	depth, start := 0, 0
	for index, char := range peers {
		switch char {
		case '(':
			if depth == 0 {
				start = index + 1
			}
			depth++
		case ')':
			depth--
			if depth < 0 {
				return "", util.ErrInvalidFormat
			}
			if depth == 0 {
				peer, err := s.maskKey(peers[start:index])
				if err != nil {
					return "", err
				}
				masked.WriteString("(" + peer + ")")
			}
		}
	}
	if depth != 0 {
		return "", util.ErrInvalidFormat
	}

	return masked.String(), nil
}

// maskReference masks a dependency's value. It points at a package key, in full before v9 and
// without the name from v9, so it has to get the same placeholder as that key.
func (s *pnpmLockSanitizer) maskReference(name, value string) (string, error) {
	base, peers := splitPnpmPeers(value)
	maskedPeers, err := s.maskPeers(peers)
	if err != nil {
		return "", err
	}

	if masked, found := s.baseKeys[name+"@"+base]; found {
		return strings.TrimPrefix(masked, name+"@") + maskedPeers, nil
	}
	if masked, found := s.baseKeys[base]; found {
		return masked + maskedPeers, nil
	}
	if pnpmSlashKeyRegex.MatchString(base) {
		return base + maskedPeers, nil
	}

	return maskPnpmVersion(base) + maskedPeers, nil
}

func maskPnpmVersion(value string) string {
	if pnpmVersionRegex.MatchString(value) {
		return value
	}

	return util.MaskValue(value)
}

func splitPnpmPeers(key string) (string, string) {
	if index := strings.Index(key, "("); index >= 0 {
		return key[:index], key[index:]
	}

	return key, ""
}
