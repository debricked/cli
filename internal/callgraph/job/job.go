package job

import error "github.com/debricked/cli/internal/io/err"

type IJob interface {
	GetFiles() []string
	GetDir() string
	Errors() error.IErrors
	Run()
	ReceiveStatus() chan string
}
