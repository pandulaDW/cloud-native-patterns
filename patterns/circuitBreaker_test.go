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

func TestCircuitBreaker(t *testing.T) {
	ctx := context.Background()

	t.Run("success response falls through the cb", func(t *testing.T) {
		mockCircuitSuccess := createMockCircuitBreaker(false)
		res, err := mockCircuitSuccess(context.Background())
		assert.Equal(t, "ok", res)
		assert.Nil(t, err)
	})

	t.Run("error response falls through the cb before reaching the failure threshold",
		func(t *testing.T) {
			mockCircuit := createMockCircuitBreaker(true)
			for i := 0; i < 5; i++ {
				res, err := mockCircuit(ctx)
				assert.NotNil(t, err)
				assert.Equal(t, "service error", err.Error())
				assert.Equal(t, "", res)
			}
		})

	t.Run("cb will break the response if a retry comes before the 2-second threshold",
		func(t *testing.T) {
			mockCircuit := createMockCircuitBreaker(true)
			for i := 0; i < 5; i++ {
				_, _ = mockCircuit(ctx)
			}

			for i := 0; i < 3; i++ {
				res, err := mockCircuit(ctx)
				assert.NotNil(t, err)
				assert.Equal(t, "service unreachable", err.Error())
				assert.Equal(t, "", res)
			}
		})

	t.Run("cb will let the program flow if a retry comes after the 2-second threshold",
		func(t *testing.T) {
			mockCircuit := createMockCircuitBreaker(true)
			for i := 0; i < 5; i++ {
				_, _ = mockCircuit(ctx)
			}

			time.Sleep(time.Second * 2)
			_, err := mockCircuit(ctx)
			assert.Equal(t, "service error", err.Error())
		})
}
