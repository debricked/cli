package cmderror

type CommandError struct {
	Code int
	Err  error
}

func (e CommandError) Error() string {
	return e.Err.Error()
}

func (e CommandError) Unwrap() error {
	return e.Err
}
