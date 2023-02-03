package maven

type Job struct {
	file       string
	cmdFactory ICmdFactory
	err        error
}

func NewJob(file string, cmdFactory ICmdFactory) *Job {
	return &Job{file: file, cmdFactory: cmdFactory}
}

func (j *Job) File() string {
	return j.file
}

func (j *Job) Error() error {
	return j.err
}

func (j *Job) Run() {
	cmd, err := j.cmdFactory.MakeDependencyTreeCmd()
	if err != nil {
		j.err = err

		return
	}
	_, err = cmd.Output()
	if err != nil {
		j.err = err

		return
	}
}
