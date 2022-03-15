package program

import (
	"context"
)

func RunFromHCL(ctx context.Context, path string) (func(), error) {
	c, err := ParseHCL(path)
	if err != nil {
		return func() {}, err
	}
	return RunInBackground(ctx, c.Programs...)
}

// RunInBackground runs a list of Programs in background goroutines.
// the returned func() can be used to stop the programs when desired.
func RunInBackground(ctx context.Context, programs ...Program) (func(), error) {
	// TODO: custom cancel func, try to SIGTERM first before KILL
	ctx, stop := context.WithCancel(ctx)

	// channels for start/check errors and program completion
	runErrCh := make(chan error, len(programs))
	checkErrCh := make(chan error, len(programs))
	doneCh := make(chan struct{}, len(programs))

	// start the programs and their checks.
	for _, p := range programs {
		go p.Start(ctx, doneCh, runErrCh)
		go p.RetryCheck(ctx, checkErrCh)
	}

	// this is our stop func, it will block until all programs are killed
	stopped := 0
	stopAndWait := func() {
		// already done
		if stopped == len(programs) {
			return
		}

		// ensure context is canceled for all programs
		stop()

		// wait for all of them to complete
		for range programs {
			<-doneCh
			stopped++
		}
	}

	// we're done starting when either anything errors,
	// or when we get one passed check per program.
	var err error
	passedChecks := 0
startLoop:
	for {
		select {
		case err = <-runErrCh:
			if err != nil {
				stopAndWait()
				break startLoop
			}
		case err = <-checkErrCh:
			if err != nil {
				stopAndWait()
				break startLoop
			}
			passedChecks++
			if passedChecks == len(programs) {
				break startLoop
			}
		}
	}

	return stopAndWait, err
}
