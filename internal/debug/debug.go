package debug

import (
	"log"
	"os"

	"github.com/fatih/color"
)

func Log(message string, debug bool) {
	if debug {
		DebugLogger := log.New(os.Stderr, "DEBUG: ", log.Ldate|log.Ltime)
		DebugLogger.Println(color.BlueString(message))
	}
}
