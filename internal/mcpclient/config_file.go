package mcpclient

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/debricked/cli/internal/io"
)

// BackupSuffix is appended to a config file's path to name its pre-modification backup.
const BackupSuffix = ".bak"

// ReadConfig reads and parses an existing client config file as a generic JSON object,
// preserving any keys this package doesn't understand. Returns existed=false (with an
// empty config) if the file doesn't exist yet. Malformed JSON is a hard error: this
// package never attempts to recover/rewrite a file it can't fully understand.
func ReadConfig(fs io.IFileSystem, path string) (config map[string]any, existed bool, err error) {
	raw, err := fs.ReadFile(path)
	if err != nil {
		if fs.IsNotExist(err) {
			return map[string]any{}, false, nil
		}

		return nil, false, fmt.Errorf("failed to read %s: %w", path, err)
	}

	config = map[string]any{}
	if len(raw) > 0 {
		if err := json.Unmarshal(raw, &config); err != nil {
			return nil, false, fmt.Errorf("%s is not valid JSON, please fix or remove it manually: %w", path, err)
		}
	}

	return config, true, nil
}

// WriteConfig marshals config as indented JSON and writes it to path, creating parent
// directories as needed. If existed is true, the previous file contents are backed up
// to path+".bak" (overwriting any previous backup) before the new content is written.
func WriteConfig(fs io.IFileSystem, path string, config map[string]any, existed bool) error {
	if existed {
		if err := backupConfig(fs, path); err != nil {
			return err
		}
	}

	if err := fs.MkdirAll(filepath.Dir(path), 0750); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", path, err)
	}

	raw, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode %s: %w", path, err)
	}
	raw = append(raw, '\n')

	if err := fs.FsWriteFile(path, raw, 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}

	return nil
}

func backupConfig(fs io.IFileSystem, path string) error {
	previous, err := fs.ReadFile(path)
	if err != nil {
		if fs.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("failed to read %s for backup: %w", path, err)
	}

	if err := fs.FsWriteFile(path+BackupSuffix, previous, 0644); err != nil {
		return fmt.Errorf("failed to back up %s: %w", path, err)
	}

	return nil
}

// HasServerEntry returns the existing entry map for name under mapKey, if any.
func HasServerEntry(config map[string]any, mapKey, name string) (map[string]any, bool) {
	servers, ok := config[mapKey].(map[string]any)
	if !ok {
		return nil, false
	}
	entry, ok := servers[name].(map[string]any)

	return entry, ok
}

// UpsertServerEntry merges fields into the named entry under mapKey, creating the
// entry and/or the servers map as needed. Fields already present on the entry that
// aren't part of fields (e.g. user-added env/cwd) are left untouched. Returns
// changed=false if every field already matched, so callers can skip an unnecessary write.
func UpsertServerEntry(config map[string]any, mapKey, name string, fields map[string]any) bool {
	servers, ok := config[mapKey].(map[string]any)
	if !ok {
		servers = map[string]any{}
	}

	entry, ok := servers[name].(map[string]any)
	if !ok {
		entry = map[string]any{}
	}

	changed := false
	for key, value := range fields {
		if existing, present := entry[key]; !present || !valuesEqual(existing, value) {
			entry[key] = value
			changed = true
		}
	}

	servers[name] = entry
	config[mapKey] = servers

	return changed
}

// RemoveServerEntry deletes the named entry under mapKey, if present. Returns
// removed=false if there was nothing to remove.
func RemoveServerEntry(config map[string]any, mapKey, name string) bool {
	servers, ok := config[mapKey].(map[string]any)
	if !ok {
		return false
	}
	if _, present := servers[name]; !present {
		return false
	}

	delete(servers, name)
	config[mapKey] = servers

	return true
}

// valuesEqual compares two JSON-compatible values by their marshaled form, since a
// value read back from JSON (e.g. []any{"a","b"}) won't be reflect.DeepEqual to the
// Go-native value we're comparing it against (e.g. []string{"a","b"}).
func valuesEqual(a, b any) bool {
	aj, aErr := json.Marshal(a)
	bj, bErr := json.Marshal(b)
	if aErr != nil || bErr != nil {
		return false
	}

	return string(aj) == string(bj)
}

// EntryCommand extracts the command/args a config entry would actually invoke.
func EntryCommand(entry map[string]any) (string, []string, error) {
	command, ok := entry["command"].(string)
	if !ok || command == "" {
		return "", nil, fmt.Errorf("missing or invalid \"command\"")
	}

	rawArgs, _ := entry["args"].([]any)
	args := make([]string, 0, len(rawArgs))
	for _, rawArg := range rawArgs {
		arg, ok := rawArg.(string)
		if !ok {
			return "", nil, fmt.Errorf("\"args\" must be a list of strings")
		}
		args = append(args, arg)
	}

	return command, args, nil
}
