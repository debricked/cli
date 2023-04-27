package job

import error "github.com/debricked/cli/pkg/io/err"

type IJob interface {
	GetFiles() []string
	Errors() error.IErrors
	Run()
	ReceiveStatus() chan string
}
