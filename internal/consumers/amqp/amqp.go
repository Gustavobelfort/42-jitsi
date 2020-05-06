package amqp

import (
	"context"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gustavobelfort/42-jitsi/internal/handler"
	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"

	"github.com/gustavobelfort/42-jitsi/internal/consumers"
)

// amqpConsumer is a little interface that narrows a channel to the usage we make of it.
//
// It helps mock it in the tests.
type amqpConsumer interface {
	Consume(queue, consumerTag string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error)
}

// Router will consume the scale teams from a rabbitmq queue by exposing a endpoint.
type AMQP struct {
	channel amqpConsumer
	queue   string
	args    amqp.Table

	timeout time.Duration

	ctx    context.Context
	cancel context.CancelFunc

	starting chan struct{}
	stopping chan struct{}

	consumerTag string

	handler handler.ScaleTeamHandler

	mu *sync.Mutex
}

var consumerSeq uint64

// NewAMQP returns a new rabbitmq consumer.
func NewAMQP(channel *amqp.Channel, queue string, args amqp.Table, hdl handler.ScaleTeamHandler, timeout time.Duration) consumers.Consumer {
	return &AMQP{
		channel: channel,
		queue:   queue,
		args:    args,

		timeout: timeout,

		starting: make(chan struct{}),
		stopping: make(chan struct{}),

		handler: hdl,

		mu: new(sync.Mutex),
	}
}

func (c *AMQP) setContext() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.ctx != nil {
		return AlreadyStartedError
	}
	c.consumerTag = "ctag-42jitsi-" + strconv.FormatUint(atomic.AddUint64(&consumerSeq, 1), 16)

	c.ctx = logging.ContextWithFields(context.Background(), logrus.Fields{
		"consumer": c.consumerTag,
		"queue":    c.queue,
	})
	c.ctx, c.cancel = context.WithCancel(c.ctx)

	return nil
}

func (c *AMQP) waitStart() {
	<-c.starting
}

// Start tries to start the consumer. If it's already started, it returns `AlreadyStartedError`.
func (c *AMQP) Start() error {
	if err := c.setContext(); err != nil {
		return err
	}
	deliveries, err := c.channel.Consume(c.queue, c.consumerTag, false, false, false, false, c.args)
	if err != nil {
		return err
	}
	close(c.starting)
	logging.ContextLog(c.ctx, logrus.StandardLogger()).Info("starting consuming")
	return c.consume(deliveries, c.stopping)
}

// Stop faithfully stops the consumer with a timeout of 20 seconds.
func (c *AMQP) Stop() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.cancel == nil {
		return NotStartedError
	}

	logging.ContextLog(c.ctx, logrus.StandardLogger()).WithField("timeout", time.Second*20).Info("shutting down consumer")
	c.cancel()
	ticker := time.NewTicker(time.Second * 20)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		return StopTimeoutError
	case <-c.stopping:
	}
	return nil
}
