package utility

import (
	"testing"
	"time"
)

// 模拟一个操作，成功于第 attemptSuccess 次调用
func mockDoFactory(successAt int) func() bool {
	attempt := 0
	return func() bool {
		attempt++
		return attempt >= successAt
	}
}

// 固定退避策略：每次都等同样的时间
func fixedBackoff(duration time.Duration) func(int) {
	return func(int) {
		time.Sleep(duration)
	}
}

func TestRetryWithBackoff_SuccessOnFirstTry(t *testing.T) {
	success := RetryWithBackoff(
		mockDoFactory(1),
		3,
		fixedBackoff(0),
	)

	if !success {
		t.Errorf("expected success on first try, got failure")
	}
}

func TestRetryWithBackoff_SuccessAfterRetries(t *testing.T) {
	success := RetryWithBackoff(
		mockDoFactory(2),
		3,
		fixedBackoff(0),
	)

	if !success {
		t.Errorf("expected success after retry, got failure")
	}
}

func TestRetryWithBackoff_FailureAfterMaxRetries(t *testing.T) {
	success := RetryWithBackoff(
		mockDoFactory(5), // Will never succeed within 3 retries
		3,
		fixedBackoff(0),
	)

	if success {
		t.Errorf("expected failure after max retries, got success")
	}
}

func TestRetryWithBackoff_ZeroRetries(t *testing.T) {
	success := RetryWithBackoff(
		mockDoFactory(1),
		0,
		fixedBackoff(0),
	)

	if success {
		t.Errorf("expected failure with zero retries, got success")
	}
}
