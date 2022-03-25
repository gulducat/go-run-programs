package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gulducat/go-run-programs/program"
)

func main() {
	os.Exit(CLI(os.Args))
}

func CLI(args []string) int {
	if len(args) < 2 {
		fmt.Println("gotta provide an hcl file ok?")
		return 1
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	stop, err := program.RunFromHCL(context.Background(), args[1])
	defer stop()

	if err != nil {
		fmt.Println("error:", err)
		return 1
	}

	s := <-sig
	log.Println("signal:", s)
	return 0
}
