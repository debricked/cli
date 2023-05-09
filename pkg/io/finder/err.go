package finder

type SetupScriptError struct {
	message string
}

type SetupWalkError struct {
	message string
}

type SetupSubprojectError struct {
	message string
}

func (e SetupScriptError) Error() string {

	return e.message
}

func (e SetupWalkError) Error() string {

	return e.message
}

func (e SetupSubprojectError) Error() string {

	return e.message
}

type SetupError []error

func (e SetupError) Error() string {
	var s string
	for _, err := range e {
		s += err.Error() + "\n"
	}

	return s
}
