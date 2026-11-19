package python

import "github.com/debricked/cli/internal/sanitize/pm/util"

func sanitizePipfileLock(content []byte, _ *util.FileContext) ([]byte, error) {
	lock, err := util.ParseJSONObject(content)
	if err != nil {
		return nil, util.ErrInvalidFormat
	}

	sanitized := util.NewOrderedObject()
	sectionsFound := false
	for _, key := range lock.Keys {
		if key != "default" && key != "develop" {
			continue
		}
		section, isObject := lock.Values[key].(*util.OrderedObject)
		if !isObject {
			return nil, util.ErrInvalidFormat
		}
		entries := util.NewOrderedObject()
		for _, name := range section.Keys {
			entry, isObject := section.Values[name].(*util.OrderedObject)
			if !isObject {
				return nil, util.ErrInvalidFormat
			}
			entries.Set(name, sanitizePipfileEntry(entry))
		}
		sanitized.Set(key, entries)
		sectionsFound = true
	}
	if !sectionsFound {
		return nil, util.ErrInvalidFormat
	}

	return util.EncodeJSON(sanitized)
}

func sanitizePipfileEntry(entry *util.OrderedObject) *util.OrderedObject {
	_, hasVersion := entry.Values["version"]
	sanitized := util.NewOrderedObject()
	for _, key := range entry.Keys {
		switch key {
		case "version":
			sanitized.Set(key, entry.Values[key])
		case "git":
			if hasVersion {
				sanitized.Set(key, "redacted")
			}
		}
	}

	return sanitized
}
