package program_test

import (
	"context"
	"fmt"
	"runtime"

	"github.com/gulducat/go-run-programs/program"
)

func Example() {
	sleep := "sleep"
	if runtime.GOOS == "windows" {
		sleep = "powershell -command Start-Sleep -Seconds"
	}

	ctx := context.Background()
	p := program.Program{
		Name:    "sleepy-server",
		Command: program.Command(sleep + " 30"),
		Check:   program.Command(sleep + " 0"),
	}

	stop, err := program.RunInBackground(ctx, p)
	defer stop()
	if err != nil {
		fmt.Println(err)
	}

	// do things with the server here.

	// Output:
}
