package main

import (
	"log"
	"os"

	"github.com/gulducat/go-run-programs/hcl"
)

func main() {
	os.Exit(CLI(os.Args))
}

func CLI(args []string) int {
	if len(args) < 2 {
		log.Println("gotta provide an hcl file ok?")
		return 1
	}

	stop, err := hcl.RunFromHCL(args[1])
	defer stop()
	if err != nil {
		return 1
	}

	// let the good times roll
	select {}
}
