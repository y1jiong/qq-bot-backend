package utility

import (
	"time"

	"github.com/gogf/gf/v2/util/grand"
)

func RetryWithBackoff(do func() bool, maxAttempts int, backoffStrategy func(int)) (success bool) {
	for i := range maxAttempts {
		if do() {
			return true
		}
		if i == maxAttempts-1 {
			break
		}
		backoffStrategy(i)
	}
	return
}

func ExponentialBackoffWithJitter(attempt int) {
	time.Sleep(1<<attempt*time.Second +
		time.Duration(grand.Intn(1000))*time.Millisecond)
}
