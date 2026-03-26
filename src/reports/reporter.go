package reporter

import (
	"fmt"
	"os"
	"strconv"
)

var HadError = false

func Report(line int, where string, message string) {
	fmt.Fprintln(os.Stderr, "[line: "+strconv.Itoa(line)+"] error"+where+": "+message)
	HadError = true
}

func Error(line int, message string) {
	Report(line, "", message)
}

func Clear() {
	HadError = false
}
