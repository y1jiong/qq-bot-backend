package utility

import (
	"github.com/gogf/gf/v2/util/grand"
	"time"
)

func RetryWithBackoff(do func() bool, maxRetries int, backoffStrategy func(int)) (success bool) {
	for i := range maxRetries {
		if do() {
			return true
		}
		if i == maxRetries-1 {
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
