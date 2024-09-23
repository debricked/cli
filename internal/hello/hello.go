package hello

import (
	"strings"
)

type IGreeter interface {
	Greeting(string) string
}

type DebrickedGreeter struct{}

func NewDebrickedGreeter() DebrickedGreeter {
	return DebrickedGreeter{}
}

func (dg DebrickedGreeter) Greeting(name string) string {
	var sb strings.Builder
	sb.WriteString("Hello ")
	sb.WriteString(name)
	sb.WriteString("!")

	return sb.String()
}
