package cmderror

type CommandError struct {
	Code int
	Err  error
}

func (e CommandError) Error() string {
	return e.Err.Error()
}
