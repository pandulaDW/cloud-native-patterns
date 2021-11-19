package patterns

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Circuit func(ctx context.Context) (string, error)

func Breaker(circuit Circuit, failureThreshold uint) Circuit {
	consecutiveFailures := 0
	lastAttempt := time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock() // Establish a read lock

		d := consecutiveFailures - int(failureThreshold)

		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}

		m.RUnlock()                   // Release read lock
		response, err := circuit(ctx) // Issue request
		m.Lock()                      // Lock around shared
		defer m.Unlock()

		lastAttempt = time.Now() // Record time of attempt
		if err != nil {          // Circuit returned an error, so we count the failure and return
			consecutiveFailures++
			return response, err
		}

		consecutiveFailures = 0 //reset failure counter

		return response, nil
	}
}
