package patterns

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func createCircuit() Circuit {
	callCount := 0
	return func(_ context.Context) (string, error) {
		callCount++
		return fmt.Sprintf("ok:%d", callCount), nil
	}
}

func TestDebounceFirst(t *testing.T) {
	ctx := context.Background()

	t.Run("initial function call will be forwarded as expected", func(t *testing.T) {
		circuit := DebounceFirst(createCircuit(), 1*time.Second)
		res, err := circuit(ctx)
		assert.Equal(t, "ok:1", res)
		assert.Nil(t, err)
	})

	t.Run("subsequent function calls within threshold will be resolved from cache", func(t *testing.T) {
		circuit := DebounceFirst(createCircuit(), 1*time.Second)
		_, _ = circuit(ctx)

		for i := 0; i < 5; i++ {
			res, err := circuit(ctx)
			assert.Equal(t, "ok:1", res)
			assert.Nil(t, err)
		}
	})

	t.Run("subsequent function calls which exceeds threshold will return called response", func(t *testing.T) {
		circuit := DebounceFirst(createCircuit(), 10*time.Millisecond)
		_, _ = circuit(ctx)

		time.Sleep(200 * time.Millisecond)
		res, err := circuit(ctx)
		assert.Equal(t, "ok:2", res)
		assert.Nil(t, err)
	})
}
