package debug

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func Print(message string, debug bool) {
	if debug {
		debugMessage := strings.Join([]string{"DEBUG: ", message, "\n"}, "")
		fmt.Fprint(os.Stderr, color.BlueString(debugMessage))
	}
}
