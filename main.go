package main

import (
	"dq/cmd"
	"fmt"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintln(err))
	}
}
