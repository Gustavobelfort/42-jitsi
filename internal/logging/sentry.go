package logging

import (
	"errors"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

// sentryHub narrows down the sentry.Hub structure methods to the methods needed by SentryHook.
type sentryHub interface {
	AddBreadcrumb(breadcrumb *sentry.Breadcrumb, hint *sentry.BreadcrumbHint)
	CaptureException(exception error) *sentry.EventID
	Flush(timeout time.Duration) bool
}

// SentryHook is a custom logrus hook that works with the latest sentry package for go.
type SentryHook struct {
	hub         sentryHub
	errorLevels map[logrus.Level]interface{}
}

func (hook *SentryHook) handleError(entry *logrus.Entry) {
	err := errors.New(entry.Message)

	if _, ok := entry.Data["error"]; ok {
		if msgErr, ok := entry.Data["error"].(error); ok && msgErr != nil {
			err = msgErr
		} else if msgErr, ok := entry.Data["error"].(string); ok && msgErr != "" {
			err = errors.New(msgErr)
		} else if msgErr, ok := entry.Data["error"].(fmt.Stringer); ok && msgErr != nil {
			err = errors.New(msgErr.String())
		}
	}

	hook.hub.CaptureException(err)
}

func (hook *SentryHook) handleBreadcrumb(entry *logrus.Entry, breadcrumbType string) {
	category := ContextGetSentryCategory(entry.Context)
	if category == "" {
		category = "default"
	}
	breadcrumb := &sentry.Breadcrumb{
		Category:  category,
		Data:      entry.Data,
		Level:     sentry.Level(entry.Level.String()),
		Message:   entry.Message,
		Timestamp: entry.Time,
		Type:      breadcrumbType,
	}
	hook.hub.AddBreadcrumb(breadcrumb, nil)
}

// Fire fire the logrus hook. It adds the current logging as a breadcrumb. And creates an event if the entry's level
// is part of the previously set error levels.
func (hook *SentryHook) Fire(entry *logrus.Entry) error {
	if _, ok := hook.errorLevels[entry.Level]; ok {
		hook.handleBreadcrumb(entry, "error")
		hook.handleError(entry)
		return nil
	}
	hook.handleBreadcrumb(entry, "default")
	return nil
}

func (hook *SentryHook) addErrorLevels(levels []logrus.Level) {
	if hook.errorLevels == nil {
		hook.errorLevels = make(map[logrus.Level]interface{})
	}
	for _, level := range levels {
		hook.errorLevels[level] = struct{}{}
	}
}

// Levels returns all the levels as all logs should be passed as breadcrumb.
func (hook *SentryHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// AddSentryHook adds a sentry hook to the given logger.
func AddSentryHook(hub *sentry.Hub, logger *logrus.Logger, levels []logrus.Level) {
	hook := &SentryHook{
		hub: hub,
	}
	hook.addErrorLevels(levels)
	logger.AddHook(hook)
	logrus.DeferExitHandler(func() { hub.Flush(time.Second * 5) })
}
