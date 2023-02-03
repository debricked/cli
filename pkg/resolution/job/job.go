package job

type IJob interface {
	File() string
	Error() error
	Run()
}
