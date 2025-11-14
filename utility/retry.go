package utility

import (
	"context"
	"errors"
	"time"

	"github.com/gogf/gf/v2/util/grand"
)

var ErrMaxRetryExceeded = errors.New("max retry exceeded")

func RetryWithBackoff(ctx context.Context, do func() bool, maxAttempts int, backoffStrategy func(int)) error {
	for i := range maxAttempts {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if do() {
			return nil
		}
		if i == maxAttempts-1 {
			break
		}
		if backoffStrategy != nil {
			backoffStrategy(i)
		}
	}
	return ErrMaxRetryExceeded
}

func ExponentialBackoffWithJitter(ctx context.Context) func(attempt int) {
	return func(attempt int) {
		select {
		case <-ctx.Done():
		case <-time.After(1<<attempt*time.Second + time.Duration(grand.Intn(1000))*time.Millisecond):
		}
	}
}
