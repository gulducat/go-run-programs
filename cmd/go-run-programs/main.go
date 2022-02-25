package main

import (
	"context"
	"log"
	"os"

	"github.com/gulducat/go-run-programs/program"
)

func main() {
	os.Exit(CLI(os.Args))
}

func CLI(args []string) int {
	if len(args) < 2 {
		log.Println("gotta provide an hcl file ok?")
		return 1
	}

	stop, err := program.RunFromHCL(context.Background(), args[1])
	defer stop()
	if err != nil {
		return 1
	}

	// let the good times roll
	select {}
}
