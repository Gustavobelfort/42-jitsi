package logging

import (
	"errors"

	"github.com/sirupsen/logrus"
)

// LogError logs an error according to its type.
//
// If the error wraps in a `WithLogError` it will use its embedded parameters to log at a specific levels
// and with specific fields.
//
// If the error doesn't wrap `WithLogError` it will log the error on `logrus.ErrorLevel` and simply add the
// error field.
func LogError(logger logrus.FieldLogger, err error, message string) {
	if err == nil {
		return
	}

	var (
		logLevel  = logrus.ErrorLevel
		logFields = make(logrus.Fields)
		logErr    = &WithLogError{}
	)

	if errors.As(err, &logErr) {
		logLevel = logErr.LogLevel
		logFields = logErr.ExtraFields
	}
	logger.WithFields(logFields).WithError(err).Logf(logLevel, "%s: %v\n", message, err)
	if logLevel == logrus.FatalLevel {
		logrus.Exit(1)
	}
}

// WithLog wraps the passed error into a `WithLogError` with the given level and fields.
//
// If the err is already a `WithLogError`, it will simply merge the contextual fields. Conflicting contextual
// fields will be overwritten.
func WithLog(err error, logLevel logrus.Level, logFields logrus.Fields) error {
	logErr := &WithLogError{}

	if err == nil {
		return nil
	}

	if logFields == nil {
		logFields = make(logrus.Fields)
	}

	if errors.As(err, &logErr) {
		for key, value := range logFields {
			logErr.ExtraFields[key] = value
		}
		return err
	}

	return &WithLogError{
		LogLevel:    logLevel,
		ExtraFields: logFields,
		err:         err,
	}
}

// WithLogError allows to store log information in the returned error.
type WithLogError struct {
	LogLevel    logrus.Level
	ExtraFields logrus.Fields

	err error
}

// Unwrap returns the underlying error.
//
// It implements the `errors.Wrapper` interface.
func (err *WithLogError) Unwrap() error {
	return err.err
}

// Returns the underlying error's string representation.
func (err *WithLogError) Error() string {
	return err.err.Error()
}
