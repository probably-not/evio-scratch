//go:build !go1.16
// +build !go1.16

package cancellation

import (
	"context"
	"os"
	"os/signal"
)

func CreateCancelContext() context.Context {
	schan := make(chan os.Signal, 1)
	ctx, cancel := context.WithCancel(context.Background())

	signal.Notify(schan, os.Interrupt)
	go func() {
		<-schan
		cancel()
	}()

	return ctx
}
