//go:build go1.16
// +build go1.16

package cancellation

import (
	"context"
	"os"
	"os/signal"
)

func CreateCancelContext() context.Context {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	return ctx
}
