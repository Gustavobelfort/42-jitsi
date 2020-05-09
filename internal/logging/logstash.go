package logging

import (
	"net"

	logrustash "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
)

// LogstashHook wraps the base logstash hook to enable level configuration.
type LogstashHook struct {
	logrus.Hook
	levels []logrus.Level
}

// Levels return the configured levels of the logstash hook wrapper.
func (hook *LogstashHook) Levels() []logrus.Level {
	return hook.levels
}

// AddLogstashHook adds a logstash hook to the given logger.
func AddLogstashHook(protocol, address string, logger *logrus.Logger, fields logrus.Fields, levels []logrus.Level) error {
	conn, err := net.Dial(protocol, address)
	if err != nil {
		return err
	}
	hook := &LogstashHook{
		Hook:   logrustash.New(conn, logrustash.DefaultFormatter(fields)),
		levels: levels,
	}
	logger.AddHook(hook)
	return nil
}
