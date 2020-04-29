package logging

import (
	"errors"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogError(t *testing.T) {

	logger, hook := test.NewNullLogger()
	expectedError := errors.New("testing error")
	expectedMessage := "message"
	expectedLogLevel := logrus.ErrorLevel
	expectedFields := logrus.Fields{
		"error": expectedError.Error(),
	}

	assertExpectations := func(t *testing.T) {
		entry := hook.LastEntry()
		assert.Equal(t, expectedLogLevel, entry.Level)
		assert.Equal(t, fmt.Sprintf("%s: %v\n", expectedMessage, expectedError), entry.Message)

		for key, expectedValue := range expectedFields {
			value, ok := entry.Data[key]
			if assert.True(t, ok) {
				assert.Equal(t, expectedValue, value)
			}
		}

	}

	t.Run("NilError", func(t *testing.T) {
		hook.Reset()

		LogError(logger, nil, expectedMessage)
		assert.Empty(t, hook.AllEntries())
	})

	t.Run("SimpleError", func(t *testing.T) {
		hook.Reset()

		LogError(logger, expectedError, expectedMessage)
		assertExpectations(t)
	})

	t.Run("WithLogError", func(t *testing.T) {
		hook.Reset()

		expectedLogLevel = logrus.WarnLevel
		expectedFields = logrus.Fields{
			"error": expectedError.Error(),
			"key":   "value",
		}
		expectedError = &WithLogError{
			LogLevel:    expectedLogLevel,
			ExtraFields: expectedFields,
			err:         expectedError,
		}

		LogError(logger, expectedError, expectedMessage)
		assertExpectations(t)
	})
}

func TestWithLogError(t *testing.T) {

	require.Implements(t, (*error)(nil), &WithLogError{})

	t.Run("NilWrap", func(t *testing.T) {
		expectedError := error(nil)
		expectedLogLevel := logrus.WarnLevel
		expectedFields := logrus.Fields{
			"key": "value",
		}

		err := WithLog(expectedError, expectedLogLevel, expectedFields)

		assert.Nil(t, err)
	})

	t.Run("NoFieldWrap", func(t *testing.T) {
		expectedError := errors.New("testing error")
		expectedLogLevel := logrus.WarnLevel
		expectedFields := logrus.Fields{}

		err := WithLog(expectedError, expectedLogLevel, nil)

		assert.Equal(t, expectedError, err.(*WithLogError).err)
		assert.Equal(t, expectedError, errors.Unwrap(err))
		assert.Equal(t, expectedError.Error(), err.Error())
		assert.Equal(t, expectedLogLevel, err.(*WithLogError).LogLevel)
		assert.Equal(t, expectedFields, err.(*WithLogError).ExtraFields)
	})

	t.Run("SimpleWrap", func(t *testing.T) {
		expectedError := errors.New("testing error")
		expectedLogLevel := logrus.WarnLevel
		expectedFields := logrus.Fields{
			"key": "value",
		}

		err := WithLog(expectedError, expectedLogLevel, expectedFields)

		assert.Equal(t, expectedError, err.(*WithLogError).err)
		assert.Equal(t, expectedError, errors.Unwrap(err))
		assert.Equal(t, expectedError.Error(), err.Error())
		assert.Equal(t, expectedLogLevel, err.(*WithLogError).LogLevel)
		assert.Equal(t, expectedFields, err.(*WithLogError).ExtraFields)
	})

	t.Run("DoubleWrap", func(t *testing.T) {
		expectedError := errors.New("testing error")
		expectedLogLevel := logrus.WarnLevel
		expectedFields := logrus.Fields{
			"key": "value",
		}

		err := WithLog(expectedError, expectedLogLevel, expectedFields)

		expectedFields = logrus.Fields{
			"key":      "newvalue",
			"otherkey": "othervalue",
		}

		err = WithLog(err, logrus.DebugLevel, expectedFields)

		assert.Equal(t, expectedError, err.(*WithLogError).err)
		assert.Equal(t, expectedError, errors.Unwrap(err))
		assert.Equal(t, expectedError.Error(), err.Error())
		assert.Equal(t, expectedLogLevel, err.(*WithLogError).LogLevel)
		assert.Equal(t, expectedFields, err.(*WithLogError).ExtraFields)
	})

}
