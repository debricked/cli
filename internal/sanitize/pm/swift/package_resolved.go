package swift

import (
	"strings"

	"github.com/debricked/cli/internal/sanitize/pm/util"
)

// maskedLocationPrefix replaces the host and path of a location. The backend only reads the
// last part of a location, and only when a pin has no identity.
const maskedLocationPrefix = "https://redacted.invalid/"

// sanitizePackageResolved keeps the name and version of each pin. The name is the identity, or
// the last part of the location if there is no identity. The version is state.version, or
// state.revision if there is no version. Branches, locations and hashes are removed.
//
// The backend only reads top-level pins, so a version 1 file, which nests them under "object",
// fails here rather than being uploaded and giving no packages.
func sanitizePackageResolved(content []byte, _ *util.FileContext) ([]byte, error) {
	resolved, err := util.ParseJSONObject(content)
	if err != nil {
		return nil, util.ErrInvalidFormat
	}
	pins, isList := resolved.Values["pins"].([]any)
	if !isList {
		return nil, util.ErrInvalidFormat
	}

	sanitizedPins := make([]any, 0, len(pins))
	for _, item := range pins {
		pin, isObject := item.(*util.OrderedObject)
		if !isObject {
			return nil, util.ErrInvalidFormat
		}
		sanitizedPin, err := sanitizePin(pin)
		if err != nil {
			return nil, err
		}
		sanitizedPins = append(sanitizedPins, sanitizedPin)
	}

	sanitized := util.NewOrderedObject()
	sanitized.Set("pins", sanitizedPins)

	return util.EncodeJSON(sanitized)
}

func sanitizePin(pin *util.OrderedObject) (*util.OrderedObject, error) {
	sanitized := util.NewOrderedObject()
	if identity := stringValue(pin, "identity"); identity != "" {
		sanitized.Set("identity", identity)
	} else if name := lastPathPart(stringValue(pin, "location")); name != "" {
		sanitized.Set("location", maskedLocationPrefix+name)
	} else {
		return nil, util.ErrInvalidFormat
	}

	sanitizedState := util.NewOrderedObject()
	if state, isObject := pin.Values["state"].(*util.OrderedObject); isObject {
		if version := stringValue(state, "version"); version != "" {
			sanitizedState.Set("version", version)
		} else if revision := stringValue(state, "revision"); revision != "" {
			sanitizedState.Set("revision", revision)
		}
	}
	sanitized.Set("state", sanitizedState)

	return sanitized, nil
}

// stringValue returns "" for a missing value or one that isn't a string, such as a null version.
func stringValue(object *util.OrderedObject, key string) string {
	value, _ := object.Values[key].(string)

	return value
}

// lastPathPart works for URLs and for scp-style git locations such as git@host:org/repo.git.
func lastPathPart(location string) string {
	location = strings.TrimRight(location, "/")

	return location[strings.LastIndexAny(location, "/:")+1:]
}
