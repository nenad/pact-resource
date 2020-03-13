package concourse

import (
	"fmt"
	"os"
)

func FailTask(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
