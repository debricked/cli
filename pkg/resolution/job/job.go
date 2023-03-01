package job

type IJob interface {
	GetFile() string
	Errors() IErrors
	Run()
	ReceiveStatus() chan string
}
