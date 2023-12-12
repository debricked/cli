package util

type PMJobError struct {
	err    string
	cmd    string
	doc    string
	status string
}

var InstalLPrivateDependencyMessage = "If this is a private dependency, please make sure that the debricked CLI has access to install it or pre-install it before running the debricked CLI."
var UnknownError = "No specific documentation for this problem yet. If you would like this message to more informative for this error, please create an issue here: https://github.com/debricked/cli/issues"

func (e PMJobError) Error() string {
	return e.err
}

func (e PMJobError) Command() string {
	if len(e.cmd) == 0 {
		return ""
	}

	return "`" + e.cmd + "`\n"
}

func (e PMJobError) Documentation() string {
	return e.doc + "\n"
}

func (e PMJobError) Status() string {
	return e.status
}

func (e *PMJobError) SetStatus(status string) {
	e.status = status
}

func (e *PMJobError) SetDocumentation(doc string) {
	e.doc = doc
}

func (e *PMJobError) SetCommand(cmd string) {
	e.cmd = cmd
}

func NewPMJobError(err string) *PMJobError {
	return &PMJobError{
		err:    err,
		cmd:    "",
		doc:    UnknownError,
		status: "",
	}
}
