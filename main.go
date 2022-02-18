package main

import (
	"os"

	"github.com/gulducat/go-run-programs/cmd"
)

// this is almost a basic supervisord :eye-flippy:
// todo: notice if a command goes away or is killed out of band

/* plain english: TODO: make this a proper package docstring
 * program(s) runs in backround (goroutine)
 * it's stopped when this program closes (context with cancel)
 * a command checks to see if the program is ready
 * > tests run here
 * program(s) stopped
 */

func main() {
	os.Exit(cmd.CLI(os.Args))
}
