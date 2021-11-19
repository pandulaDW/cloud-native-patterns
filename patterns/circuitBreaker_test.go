package patterns

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func createMockCircuitBreaker(isFailure bool) Circuit {
	f := func(_ context.Context) (string, error) {
		if isFailure {
			return "", errors.New("service error")
		}
		return "ok", nil
	}
	return Breaker(f, 5)
}

// Assert that success response falls through the circuit breaker
func TestCircuitBreakerSuccess(t *testing.T) {
	mockCircuitSuccess := createMockCircuitBreaker(false)
	res, err := mockCircuitSuccess(context.Background())

	assert.Equal(t, "ok", res)
	assert.Nil(t, err)
}

// Assert that error response falls through the circuit breaker before reaching the failure threshold
func TestCircuitBreakerErrorResponse(t *testing.T) {
	mockCircuit := createMockCircuitBreaker(true)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		res, err := mockCircuit(ctx)
		assert.NotNil(t, err)
		assert.Equal(t, "service error", err.Error())
		assert.Equal(t, "", res)
	}
}

// Assert that circuit breaker will break the response if a retry comes before the 2-second threshold
func TestCircuitBreakerBreakerResponse(t *testing.T) {
	mockCircuit := createMockCircuitBreaker(true)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_, _ = mockCircuit(ctx)
	}

	for i := 0; i < 3; i++ {
		res, err := mockCircuit(ctx)
		assert.NotNil(t, err)
		assert.Equal(t, "service unreachable", err.Error())
		assert.Equal(t, "", res)
	}
}

// Assert that circuit breaker will let the program flow if a retry comes after the 2-second threshold
func TestCircuitBreakerResponseAfterTimeout(t *testing.T) {
	mockCircuit := createMockCircuitBreaker(true)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		_, _ = mockCircuit(ctx)
	}

	time.Sleep(time.Second * 2)
	mockCircuit = createMockCircuitBreaker(false)
	res, err := mockCircuit(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "ok", res)
}
