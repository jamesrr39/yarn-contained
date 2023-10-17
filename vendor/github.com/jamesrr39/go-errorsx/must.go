package errorsx

import (
	"fmt"
	"os"
)

// ExitIfErr prints an error message and stack trace, and then exits the application if an error is passed in
func ExitIfErr(err Error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\nStack trace:\n%s\n", err.Error(), err.Stack())
		os.Exit(1)
	}
}
