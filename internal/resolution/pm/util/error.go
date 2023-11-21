package util

type PMJobError struct {
	err    string
	cmd    string
	doc    string
	status string
}

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
		doc:    "No specific documentation for this problem yet, please report it to us! :)",
		status: "",
	}
}
