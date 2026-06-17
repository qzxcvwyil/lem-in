package main

import (
	"fmt"
	"os"

	"01.tomorrow-school.ai/git/msabyrga/lem-in.git/internal/app"
)

func main() {
	args := os.Args[1:]

	if err := app.Run(args); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: invalid data format, %v\n", err)
		os.Exit(1)
	}
}
