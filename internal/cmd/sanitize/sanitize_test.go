package sanitize

import (
	"bytes"
	"errors"
	"testing"

	"github.com/debricked/cli/internal/cmd/cmderror"
	"github.com/debricked/cli/internal/file"
	"github.com/debricked/cli/internal/sanitize"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

type sanitizerMock struct {
	result  sanitize.Result
	err     error
	options sanitize.Options
}

func (m *sanitizerMock) Sanitize(options sanitize.Options) (sanitize.Result, error) {
	m.options = options

	return m.result, m.err
}

func execute(mock *sanitizerMock, args ...string) (string, string, error) {
	viper.Reset()
	cmd := NewSanitizeCmd(mock)
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs(args)
	err := cmd.Execute()

	return out.String(), errOut.String(), err
}

var sanitizedGroups = []file.Group{
	*file.NewGroup(
		".debricked/sanitized/web/package.json",
		nil,
		[]string{".debricked/sanitized/web/yarn.lock"},
	),
	*file.NewGroup("", nil, []string{".debricked/sanitized/net/paket.lock"}),
}

func TestNewSanitizeCmd(t *testing.T) {
	cmd := NewSanitizeCmd(&sanitizerMock{})

	assert.NotNil(t, cmd.Flags().Lookup(ExclusionFlag))
	assert.NotNil(t, cmd.Flags().Lookup(InclusionFlag))
	assert.Equal(t, "j", cmd.Flags().Lookup(JsonFlag).Shorthand)
}

func TestRunEPrintsSanitizedFilesLikeFind(t *testing.T) {
	mock := &sanitizerMock{result: sanitize.Result{Groups: sanitizedGroups}}

	out, _, err := execute(mock, "x")

	assert.NoError(t, err)
	assert.Equal(t, `.debricked/sanitized/web/package.json
 * .debricked/sanitized/web/yarn.lock
.debricked/sanitized/net/paket.lock
`, out)
}

func TestRunEPrintsSanitizedFilesAsJson(t *testing.T) {
	mock := &sanitizerMock{result: sanitize.Result{Groups: sanitizedGroups}}

	out, _, err := execute(mock, "x", "--json")

	assert.NoError(t, err)
	assert.JSONEq(t, `[
		{"manifestFile": ".debricked/sanitized/web/package.json", "lockFiles": [".debricked/sanitized/web/yarn.lock"]},
		{"manifestFile": "", "lockFiles": [".debricked/sanitized/net/paket.lock"]}
	]`, out)
}

func TestRunEPrintsSummary(t *testing.T) {
	mock := &sanitizerMock{result: sanitize.Result{
		OutputDir: "x/.debricked/sanitized",
		MapPath:   "x/.debricked/sanitize-map.json",
		Sanitized: []string{"x/yarn.lock", "x/package.json"},
		Skipped:   []sanitize.FileIssue{{Path: "x/pom.xml", Reason: "format is not supported by sanitize"}},
	}}

	out, errOut, err := execute(mock, "x")

	assert.NoError(t, err)
	assert.Empty(t, out)
	assert.Equal(t, "x", mock.options.RootPath)
	assert.Contains(t, errOut, "Sanitized: 2, skipped: 1, failed: 0")
	assert.Contains(t, errOut, "x/pom.xml: format is not supported by sanitize")
	assert.Contains(t, errOut, "x/.debricked/sanitized")
	assert.Contains(t, errOut, "x/.debricked/sanitize-map.json")
}

func TestRunEExitsWithErrorWhenAFileFailed(t *testing.T) {
	mock := &sanitizerMock{result: sanitize.Result{
		Failed: []sanitize.FileIssue{{Path: "x/yarn.lock", Reason: "content is not valid for this file type"}},
	}}

	_, errOut, err := execute(mock, "x")

	var cmdErr cmderror.CommandError
	assert.ErrorAs(t, err, &cmdErr)
	assert.Equal(t, FailureExitCode, cmdErr.Code)
	assert.Contains(t, errOut, "Sanitized: 0, skipped: 0, failed: 1")
	assert.Contains(t, errOut, "x/yarn.lock")
}

func TestRunEReturnsRunError(t *testing.T) {
	failure := errors.New("failed")

	_, _, err := execute(&sanitizerMock{err: failure})

	assert.ErrorIs(t, err, failure)
}
