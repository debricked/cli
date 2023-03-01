package job

type BaseJob struct {
	File   string
	Err    IJobError
	Status chan string
}

func (j *BaseJob) GetFile() string {
	return j.File
}

func (j *BaseJob) Error() IJobError {
	return j.Err
}

func (j *BaseJob) ReceiveStatus() chan string {
	return j.Status
}

func (j *BaseJob) SendStatus(status string) {
	j.Status <- status
}
