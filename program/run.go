package program

import "context"

// type StartCheckable interface { // todo: hmm...
// 	Start(context.Context, chan<- error)
// 	RetryCheck(chan<- error)
// }

// RunInBackground runs a list of Programs in background goroutines.
// the returned func() can be used to stop the programs when desired.
func RunInBackground(programs ...Program) (func(), error) {
	ctx, stop := context.WithCancel(context.Background())
	// an error channel lets us know if a program's main Command fails,
	// or if a Check never succeeds after retries.
	errCh := make(chan error)

	// start the programs and their checks.
	for _, p := range programs {
		go p.Start(ctx, errCh)
		go p.RetryCheck(errCh)
	}

	// wait for them to be up and ready.
	// since the program's primary Command should block indefinitely,
	// we should get 1 nil per program, from its Check.
	// e.g. with 3 programs, 3 nil's in the error channel means 3 successul Checks
	for range programs {
		// all programs must run, a single error stops everything
		if err := <-errCh; err != nil {
			stop()
			return stop, err
		}
	}

	return stop, nil
}
