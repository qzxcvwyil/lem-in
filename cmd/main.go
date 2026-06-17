package main

import (
	"log"
	"os"

	"01.tomorrow-school.ai/git/msabyrga/lem-in.git/internal/app"
)

func main() {
	args := os.Args[1:]

	if err := app.Run(args); err != nil {
		log.Fatal(err)
	}
}
