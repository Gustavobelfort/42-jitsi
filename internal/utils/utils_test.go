package utils

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWrapContext(t *testing.T) {
	t.Run("ReturnNil", func(t *testing.T) {
		err := WrapContext(context.Background(), func() error {
			return nil
		})
		assert.NoError(t, err)
	})

	t.Run("ReturnFunctionError", func(t *testing.T) {
		expectedError := errors.New("testing")

		err := WrapContext(context.Background(), func() error {
			return expectedError
		})
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("ReturnContextCanceled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())

		cancel()
		err := WrapContext(ctx, func() error {
			time.Sleep(time.Second * 2)
			return nil
		})
		assert.Error(t, err)
		assert.True(t, errors.Is(err, context.Canceled))
	})
}
