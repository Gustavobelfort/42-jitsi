package logging

import (
	"context"

	"github.com/sirupsen/logrus"
)

// This create a private key-space in the Context, meaning that only this package can get or set "contextKey" types
type contextKey int

const (
	logFieldsKey contextKey = iota
	sentryCategoryKey
)

// ContextLog returns an entry with the fields of the passed context.
func ContextLog(ctx context.Context, logger logrus.FieldLogger) *logrus.Entry {
	return logger.WithFields(ContextGetFields(ctx))
}

// ContextGetFields returns the logging fields of the context.
// If the context is nil or there are no fields, empty fields are returned.
func ContextGetFields(ctx context.Context) logrus.Fields {
	if ctx != nil {
		if fields, _ := ctx.Value(logFieldsKey).(logrus.Fields); fields != nil {
			return fields
		}
	}
	return logrus.Fields{}
}

// ContextWithFields adds logging fields to the given context.
func ContextWithFields(ctx context.Context, fields logrus.Fields) context.Context {
	oldFields := ContextGetFields(ctx)
	for key, value := range fields {
		oldFields[key] = value
	}
	return context.WithValue(ctx, logFieldsKey, oldFields)
}

// ContextWithField adds a logging field to the given context.
func ContextWithField(ctx context.Context, key string, value interface{}) context.Context {
	return ContextWithFields(ctx, logrus.Fields{key: value})
}

// ContextWithSentryCategory sets the sentry category for the given context.
func ContextWithSentryCategory(ctx context.Context, category string) context.Context {
	return context.WithValue(ctx, sentryCategoryKey, category)
}

// ContextGetSentryCategory gets the sentry category of the given context.
func ContextGetSentryCategory(ctx context.Context) string {
	category := ""
	if ctx != nil {
		category, _ = ctx.Value(sentryCategoryKey).(string)
	}
	return category
}
