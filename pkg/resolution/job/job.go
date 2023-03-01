package job

type IJob interface {
	GetFile() string
	Error() IJobError
	Run()
	ReceiveStatus() chan string
}
