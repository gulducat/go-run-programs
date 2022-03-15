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
// A map of EnvVars may be provided which will be set before Command starts.
type Program struct {
	Name         string  `hcl:"name,label"`
	Command      Command `hcl:"command"`
	Check        Command `hcl:"check"`
	CheckSeconds int     `hcl:"seconds,optional"`

	EnvVars map[string]string `hcl:"env,optional"`
}

func (p Program) String() string {
	return fmt.Sprintf("%s (%s)", p.Name, p.Command)
}

func (p Program) Start(ctx context.Context, doneCh chan<- struct{}, errCh chan<- error) {
	log.Println("starting", p)

	// hmm, do this or let the user set it in their environment...?
	for k, v := range p.EnvVars {
		if err := os.Setenv(k, v); err != nil {
			errCh <- err
		}
	}

	// happy path is this blocking forever until context is canceled
	err := p.Command.Run(ctx)
	if err != nil && err != context.Canceled {
		log.Println("errored", p, err)
		errCh <- fmt.Errorf("command '%s' stopped: %w", p.Command, err)
	}
	log.Println("done", p)
	doneCh <- struct{}{}
}

func (p Program) RetryCheck(ctx context.Context, errCh chan<- error) {
	// retry for up to ~20 seconds (default)
	sleep := time.Second
	tries := p.CheckSeconds
	if tries == 0 {
		tries = 20
	}

	var lastErr error
	for i := 0; i < tries; i++ {
		// if context becomes done, stop retrying.
		select {
		case <-ctx.Done():
			errCh <- ctx.Err()
			return
		default: // proceed
		}

		// sleep first since we start RetryCheck()ing immediately after Start()ing,
		// and it takes >0 time for the programs to become ready.
		time.Sleep(sleep)

		log.Println("checking", p, p.Check)

		// don't let the check command run for more than a few seconds
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		lastErr = p.Check.Run(ctx)
		cancel() // "the cancel function returned by context.WithTimeout should be called... (lostcancel)"
		// no error, the check was successful, so stop retrying, we're done!
		if lastErr == nil {
			break
		}
	}

	if lastErr != nil {
		log.Println("errored", p, lastErr)
	} else {
		log.Println("started", p)
	}

	errCh <- lastErr
}

// Command is a system command to run.
type Command string

func (c Command) Run(ctx context.Context) error {
	parts := strings.Split(string(c), " ")
	name := parts[0]
	args := parts[1:]

	bts, err := exec.CommandContext(ctx, name, args...).CombinedOutput()
	// if context has stopped, 'err' will be confusingly non-nil,
	// so check context first.
	if ctx.Err() != nil {
		// pass context "error"s through unwrapped
		return ctx.Err()
	}
	if err != nil {
		// exec's CombinedOutput() err is often quite vague if the command exits with error.
		// the actual output of the program, both stdout and stderr, goes to bts []byte,
		// so include that in the error message.
		return fmt.Errorf("command '%s' error: %w; output:\n%s", c, err, string(bts))
	}

	return nil
}
