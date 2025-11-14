package utility

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestRetryWithBackoff_SuccessFirstTry(t *testing.T) {
	ctx := context.Background()
	calls := 0
	backoffCalls := 0

	err := RetryWithBackoff(ctx, func() bool {
		calls++
		return true
	}, 3, func(int) {
		backoffCalls++
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
	if backoffCalls != 0 {
		t.Fatalf("expected 0 backoff calls, got %d", backoffCalls)
	}
}

func TestRetryWithBackoff_SuccessAfterRetries(t *testing.T) {
	ctx := context.Background()
	calls := 0
	backoffCalls := 0

	err := RetryWithBackoff(ctx, func() bool {
		calls++
		return calls == 3
	}, 5, func(int) {
		backoffCalls++
	})

	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
	// when success on third attempt, backoff should have been called for attempts 0 and 1
	if backoffCalls != 2 {
		t.Fatalf("expected 2 backoff calls, got %d", backoffCalls)
	}
}

func TestRetryWithBackoff_MaxRetryExceeded(t *testing.T) {
	ctx := context.Background()
	calls := 0
	backoffCalls := 0

	err := RetryWithBackoff(ctx, func() bool {
		calls++
		return false
	}, 3, func(int) {
		backoffCalls++
	})

	if !errors.Is(err, ErrMaxRetryExceeded) {
		t.Fatalf("expected ErrMaxRetryExceeded, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
	// for 3 attempts, backoff should have been called after attempt 0 and 1
	if backoffCalls != 2 {
		t.Fatalf("expected 2 backoff calls, got %d", backoffCalls)
	}
}

func TestRetryWithBackoff_ContextCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	calls := 0
	backoffCalls := 0

	// cancel before any attempt
	cancel()

	err := RetryWithBackoff(ctx, func() bool {
		calls++
		return false
	}, 5, func(int) {
		backoffCalls++
	})

	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls != 0 {
		t.Fatalf("expected 0 calls, got %d", calls)
	}
	if backoffCalls != 0 {
		t.Fatalf("expected 0 backoff calls, got %d", backoffCalls)
	}
}

func TestExponentialBackoffWithJitter_RespectsContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	backoff := ExponentialBackoffWithJitter(ctx)

	// cancel context and ensure backoff returns quickly
	cancel()

	start := time.Now()
	backoff(5)
	elapsed := time.Since(start)

	if elapsed > 50*time.Millisecond {
		t.Fatalf("expected backoff to return quickly after cancel, took %v", elapsed)
	}
}
