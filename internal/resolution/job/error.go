package job

type IError interface {
	Error() string
	Command() string
	Documentation() string
	Status() string
	SetStatus(string)
	SetDocumentation(string)
	SetCommand(string)
}

type BaseJobError struct {
	err           string
	command       string
	documentation string
	status        string
}

func (e BaseJobError) Error() string {
	return e.err
}

func (e BaseJobError) Command() string {
	return e.command
}

func (e BaseJobError) Documentation() string {
	return e.documentation + "\n"
}

func (e BaseJobError) Status() string {
	return e.status
}

func (e *BaseJobError) SetStatus(status string) {
	e.status = status
}

func (e *BaseJobError) SetDocumentation(doc string) {
	e.documentation = doc
}

func (e *BaseJobError) SetCommand(command string) {
	e.command = command
}

func NewBaseJobError(err string) *BaseJobError {
	return &BaseJobError{
		err:           err,
		command:       "",
		documentation: "",
		status:        "",
	}
}
