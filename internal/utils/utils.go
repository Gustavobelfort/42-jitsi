package utils

import (
	"context"
)

// Color const for the terminal
const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[1;33m"
)

// WrapContext is a util meant to wrap functions that are without context when the context is required.
//
// If the function returns its error before the context is canceled, then the function's error is returned.
// If the context is canceled before the enf of `f` then the context's error is returned.
func WrapContext(ctx context.Context, f func() error) error {
	var err error
	chErr := make(chan error)

	go func() {
		defer close(chErr)
		chErr <- f()
	}()

	select {
	case err = <-chErr:
		break
	case <-ctx.Done():
		err = ctx.Err()
		break
	}
	return err
}
