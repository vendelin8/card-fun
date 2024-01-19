package util

import (
	"context"
	"time"
)

const ctxTmt = time.Second * 15

// CtxTimeout creates a new context with a constant timeout.
func CtxTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), ctxTmt)
}
