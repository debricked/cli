package job

type IErrors interface {
	Warning(err IError)
	Critical(err IError)
	GetWarningErrors() []IError
	GetCriticalErrors() []IError
	GetAll() []IError
	HasError() bool
}

type Errors struct {
	title        string
	warningErrs  []IError
	criticalErrs []IError
}

func NewErrors(title string) *Errors {
	return &Errors{
		title:        title,
		warningErrs:  []IError{},
		criticalErrs: []IError{},
	}
}

func (errors *Errors) Warning(err IError) {
	errors.warningErrs = append(errors.warningErrs, err)
}

func (errors *Errors) Critical(err IError) {
	errors.criticalErrs = append(errors.criticalErrs, err)
}

func (errors *Errors) GetWarningErrors() []IError {
	return errors.warningErrs
}

func (errors *Errors) GetCriticalErrors() []IError {
	return errors.criticalErrs
}

func (errors *Errors) GetAll() []IError {
	return append(errors.warningErrs, errors.criticalErrs...)
}

func (errors *Errors) HasError() bool {
	return len(errors.criticalErrs) > 0 || len(errors.warningErrs) > 0
}
