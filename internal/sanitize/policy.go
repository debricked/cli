package sanitize

import "github.com/debricked/cli/internal/sanitize/pm/util"

// sanitizePolicyFile copies debricked_policy.json with its content unchanged. The file holds the
// customer's own policy rules, which the backend needs as they are, and no dependency data. It is
// still parsed, so a file that isn't valid JSON fails instead of being copied as it is.
func sanitizePolicyFile(content []byte, _ *util.FileContext) ([]byte, error) {
	policy, err := util.ParseJSONObject(content)
	if err != nil {
		return nil, util.ErrInvalidFormat
	}

	return util.EncodeJSON(policy)
}
