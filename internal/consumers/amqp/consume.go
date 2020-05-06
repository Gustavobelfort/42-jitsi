package amqp

import (
	"context"
	"errors"

	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

func messageContext(ctx context.Context, msg amqp.Delivery) context.Context {
	return logging.ContextWithFields(ctx, logrus.Fields{
		"model":       msg.Headers["X-Model"],
		"event":       msg.Headers["X-Event"],
		"delivery_id": msg.Headers["X-Delivery"],
	})
}

func validateMessage(msg amqp.Delivery) error {

	if model, ok := msg.Headers["X-Model"]; !ok || model != "scale_team" {
		return errors.New("unhandled model")
	}

	event, ok := msg.Headers["X-Event"]
	if !ok {
		return errors.New("unhandled event")
	}
	switch event {
	case "create", "update", "destroy":
		break
	default:
		return errors.New("unhandled event")
	}

	return nil
}

func (c *AMQP) treatMessage(ctx context.Context, msg amqp.Delivery) error {
	var err error
	ctx, cancel := context.WithTimeout(ctx, c.timeout)
	defer cancel()

	if err = validateMessage(msg); err != nil {
		return logging.WithLog(err, logrus.WarnLevel, nil)
	}

	switch msg.Headers["X-Event"].(string) {
	case "create":
		err = c.handler.HandleCreate(ctx, msg.Body)
	case "update":
		err = c.handler.HandleUpdate(ctx, msg.Body)
	case "destroy":
		err = c.handler.HandleDestroy(ctx, msg.Body)
	default:
		err = logging.WithLog(errors.New("unknown event"), logrus.WarnLevel, nil)
	}
	return err
}

func (c *AMQP) consume(deliveries <-chan amqp.Delivery, stopping chan struct{}) error {
	var err error

	defer close(stopping)

	for err == nil {
		select {
		case msg, ok := <-deliveries:
			if !ok {
				err = DeliveryChannelClosedError
				break
			}
			ctx := messageContext(c.ctx, msg)
			ctxlogger := logging.ContextLog(ctx, logrus.StandardLogger())
			ctxlogger.Info("received message")
			if err := c.treatMessage(ctx, msg); err != nil {
				logging.LogError(ctxlogger, err, "while treating the message")
				logging.LogError(ctxlogger, handleError(msg, err), "while rejecting the message")
				break
			}
			ctxlogger.Info("acknowledging message")
			logging.LogError(ctxlogger, msg.Ack(false), "while acknowledging the message")
		case <-c.ctx.Done():
			err = ConsumerStoppedError
		}
	}
	return err
}
