package amqp

import (
	"errors"

	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	// AlreadyStartedError is returned when you try to start a consumer that was already started.
	AlreadyStartedError = errors.New("the consumer was already started")
	// NotStartedError is returned when you try to stop a consumer that has not started yet.
	NotStartedError = errors.New("the consumer has not started yet")
	// ConsumerStoppedError is returned when a consumer is stopped manually.
	ConsumerStoppedError = errors.New("the consumer was stopped")
	// DeliveryChannelClosedError is returned when a consumer stopped because the server.
	DeliveryChannelClosedError = errors.New("the delivery channel was closed for some reason")
	// StopTimeoutError is returned when the consumer took more than 20 seconds to stop.
	StopTimeoutError = errors.New("the consumer timed out while stopping")
)

func handleError(msg amqp.Delivery, err error) error {
	logError := &logging.WithLogError{}
	if errors.As(err, &logError) {
		if logError.LogLevel <= logrus.WarnLevel {
			return msg.Ack(false)
		}
	}
	return msg.Reject(false)
}
