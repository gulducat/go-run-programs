package program_test

import (
	"context"
	"fmt"

	"github.com/gulducat/go-run-programs/program"
)

func Example() {
	ctx := context.Background()
	p := program.Program{
		Name:    "sleepy-server",
		Command: "sleep 30",
		Check:   "echo you are getting very very sleepy",
	}
	stop, err := program.RunInBackground(ctx, p)
	defer stop()
	if err != nil {
		fmt.Println(err)
	}
	// Output:
}
