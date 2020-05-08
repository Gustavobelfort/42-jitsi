package logging

import (
	"github.com/getsentry/sentry-go"
	"github.com/gustavobelfort/42-jitsi/internal/config"
	"github.com/sirupsen/logrus"
)

// InitiateSentry will initiate sentry logging on the standard logger based on the app's configuration.
func InitiateSentry() error {
	if !config.Conf.Sentry.Enabled {
		logrus.Debug("sentry is disabled by config")
		return nil
	}
	logrus.Info("initiating sentry")
	options := sentry.ClientOptions{
		Dsn:              config.Conf.Sentry.DSN,
		AttachStacktrace: true,
		Environment:      config.Conf.Environment,
		MaxBreadcrumbs:   100,
	}
	if err := sentry.Init(options); err != nil {
		return err
	}
	AddSentryHook(sentry.CurrentHub(), logrus.StandardLogger(), config.Conf.Sentry.Levels)
	return nil
}

// Initiate will initiate logging based on the app's configuration.
func Initiate() {
	logrus.SetLevel(config.Conf.LogLevel)
	LogError(logrus.StandardLogger(), InitiateSentry(), "while initiating sentry")
}
