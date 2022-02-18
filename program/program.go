package program

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Program's Command should run indefinitely, as with a server.
// Its Check is run several times to determine when the program is ready,
// for approximately CheckSeconds number of seconds before failing (default 20).
type Program struct {
	Name         string  `hcl:"name,label"`
	Command      Command `hcl:"command"`
	Check        Command `hcl:"check"`
	CheckSeconds int     `hcl:"seconds,optional"`

	EnvVars map[string]string `hcl:"env,optional"`
}

func (p Program) String() string {
	return p.Name
}

func (p Program) Start(ctx context.Context, errCh chan<- error) {
	// hmm, do this or let the user set it in their environment...?
	for k, v := range p.EnvVars {
		if err := os.Setenv(k, v); err != nil {
			errCh <- err
		}
	}

	// happy path is this blocking forever until context is canceled
	err := p.Command.Run(ctx)
	// so if this happens, something has gone wrong
	errCh <- fmt.Errorf("command %s stopped: %w", p.Command, err)
}

func (p Program) RetryCheck(errCh chan<- error) {
	var err error

	// retry for up to ~20 seconds (default)
	sleep := time.Second
	tries := p.CheckSeconds
	if tries == 0 {
		tries = 20
	}

	for i := 0; i < tries; i++ {
		// sleep first since we start Wait()ing immediately after Start()ing,
		// and it takes >0 time for the programs to become ready.
		time.Sleep(sleep)

		// don't let the check command run for more than a couple seconds
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		err = p.Check.Run(ctx)
		cancel() // "the cancel function returned by context.WithTimeout should be called... (lostcancel)"
		// no error, the check was successful, so stop retrying, we're done!
		if err == nil {
			errCh <- nil
			return
		}
	}

	// err will be the last Check failure
	errCh <- err
}

// Command is a system command to run.
type Command string

func (c Command) Run(ctx context.Context) error {
	log.Println("running:", c) // TODO: remove log's ?

	parts := strings.Split(string(c), " ")
	name := parts[0]
	args := parts[1:]

	bts, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	if err != nil {
		// exec's CombinedOutput() err is often quite vague if the command exits with error.
		// the actual output of the program, both stdout and stderr, goes to bts []byte,
		// so include that in the error message.
		return fmt.Errorf("command %s error: %w; output:\n%s", c, err, string(bts))
	}

	log.Println("complete:", c)

	return nil
}
