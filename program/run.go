package program

import (
	"context"
)

func RunFromHCL(ctx context.Context, path string) (func(), error) {
	c, err := ParseHCL(path)
	if err != nil {
		return func() {}, err
	}
	return RunInBackground(context.Background(), c.Programs...)
}

// RunInBackground runs a list of Programs in background goroutines.
// the returned func() can be used to stop the programs when desired.
func RunInBackground(ctx context.Context, programs ...Program) (func(), error) {
	// TODO: custom cancel func, try to SIGTERM first before KILL
	ctx, stop := context.WithCancel(ctx)

	// channels for errors and check done-ness
	errCh := make(chan error)
	doneCh := make(chan struct{})

	// start the programs and their checks.
	for _, p := range programs {
		go p.Start(ctx, errCh)
		go p.RetryCheck(ctx, doneCh, errCh)
	}

	// wait for them to be done or any error
	doneCount := 0
	for {
		select {
		// any error is bad news, stop everything
		case err := <-errCh:
			stop()
			return stop, err

		// we expect 1 done per program
		case <-doneCh:
			doneCount++
			if doneCount == len(programs) {
				return stop, nil
			}

		// any other way the context gets done
		case <-ctx.Done():
			return stop, ctx.Err()
		}
	}
}
