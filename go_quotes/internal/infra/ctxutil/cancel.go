package ctxutil

import (
	"context"
	"errors"
	"net/url"
)

func IsCancel(err error, ctx context.Context) bool {
	if err == nil {
		return false
	}
	// Parent context cancel or timeout
	if ctx != nil && ctx.Err() != nil {
		return true
	}

	// Error is canceled or deadline exceeded
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return true
	}

	// Error of HTTP *url.Error
	var ue *url.Error
	if errors.As(err, &ue) && (errors.Is(ue.Err, context.Canceled) || errors.Is(ue.Err, context.DeadlineExceeded)) {
		return true
	}
	return false
}
