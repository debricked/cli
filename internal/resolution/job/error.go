package job

type IError interface {
	Error() string
	Command() string
	Documentation() string
	Status() string
	IsCritical() bool
	SetStatus(string)
	SetDocumentation(string)
	SetCommand(string)
	SetIsCritical(bool)
}

type BaseJobError struct {
	err           string
	command       string
	documentation string
	status        string
	isCritical    bool
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

func (e BaseJobError) IsCritical() bool {
	return e.isCritical
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

func (e *BaseJobError) SetIsCritical(isCritical bool) {
	e.isCritical = isCritical
}

func NewBaseJobError(err string) *BaseJobError {
	return &BaseJobError{
		err:           err,
		command:       "",
		documentation: "",
		status:        "",
		isCritical:    true,
	}
}
