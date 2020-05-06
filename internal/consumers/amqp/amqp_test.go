package amqp

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlerMock struct {
	mock.Mock
}

func (m *HandlerMock) HandleCreate(ctx context.Context, data []byte) error {
	return m.Called(ctx, data).Error(0)
}

func (m *HandlerMock) HandleUpdate(ctx context.Context, data []byte) error {
	return m.Called(ctx, data).Error(0)
}

func (m *HandlerMock) HandleDestroy(ctx context.Context, data []byte) error {
	return m.Called(ctx, data).Error(0)
}

type ChannelMock struct {
	confirm chan struct{}
	mock.Mock
}

func (m *ChannelMock) Confirm() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	select {
	case m.confirm <- struct{}{}:
	case <-ticker.C:
	}
}

func (m *ChannelMock) Wait() {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	select {
	case <-m.confirm:
	case <-ticker.C:
	}
}

func (m *ChannelMock) Consume(queue, consumerTag string, autoAck, exclusive, noLocal, noWait bool, args amqp.Table) (<-chan amqp.Delivery, error) {
	toReturn := m.Called(queue, consumerTag, autoAck, exclusive, noLocal, noWait, args)
	return toReturn.Get(0).(chan amqp.Delivery), toReturn.Error(1)
}

func (m *ChannelMock) Ack(tag uint64, multiple bool) error {
	defer m.Confirm()
	return m.Called(tag, multiple).Error(0)
}

func (m *ChannelMock) Nack(tag uint64, multiple bool, requeue bool) error {
	defer m.Confirm()
	return m.Called(tag, multiple, requeue).Error(0)
}

func (m *ChannelMock) Reject(tag uint64, requeue bool) error {
	defer m.Confirm()
	return m.Called(tag, requeue).Error(0)
}

func TestAMQP(t *testing.T) {
	t.Run("setContext", func(t *testing.T) {
		expected := logrus.Fields{
			"queue": "test",
		}

		consumer := &AMQP{queue: expected["queue"].(string), mu: new(sync.Mutex)}
		assert.NoError(t, consumer.setContext())
		expected["consumer"] = consumer.consumerTag

		require.NotNil(t, consumer.ctx)
		assert.NotNil(t, consumer.cancel)
		assert.Equal(t, expected, logging.ContextGetFields(consumer.ctx))

		assert.Equal(t, AlreadyStartedError, consumer.setContext())
	})

	t.Run("closeDeliveries", func(t *testing.T) {
		cMock := &ChannelMock{confirm: make(chan struct{})}
		hMock := &HandlerMock{}
		defer cMock.AssertExpectations(t)
		defer hMock.AssertExpectations(t)

		expectedArgs := amqp.Table{"you": "should watch the mandalorian"}
		consumer := NewAMQP(nil, "queue", expectedArgs, hMock, time.Second*10).(*AMQP)

		consumer.channel = cMock

		deliveries := make(chan amqp.Delivery)
		cMock.On("Consume", consumer.queue, mock.Anything, false, false, false, false, consumer.args).Return(deliveries, nil).Once()

		c := make(chan error)
		go func() { c <- consumer.Start() }() // Prevent being stuck.

		consumer.waitStart()
		close(deliveries)

		ticker := time.NewTicker(time.Second * 2)
		defer ticker.Stop()

		var err error
		select {
		case err = <-c:
		case <-ticker.C:
		}
		assert.Equal(t, DeliveryChannelClosedError, err)

	})

	suite.Run(t, new(TestAMQPSuite))
}

type TestAMQPSuite struct {
	suite.Suite

	cMock *ChannelMock
	hMock *HandlerMock

	deliveries chan amqp.Delivery

	amqp *AMQP

	stop chan error
}

func (s *TestAMQPSuite) SetupSuite() {
	s.cMock = &ChannelMock{confirm: make(chan struct{})}
	s.hMock = &HandlerMock{}

	expectedArgs := amqp.Table{"you": "should watch the mandalorian"}
	s.amqp = NewAMQP(nil, "queue", expectedArgs, s.hMock, time.Second*10).(*AMQP)

	s.Nil(s.amqp.channel)
	s.Equal("queue", s.amqp.queue)
	s.Equal(expectedArgs, s.amqp.args)
	s.Require().NotNil(s.amqp.mu)
	s.Require().NotNil(s.amqp.starting)
	s.Require().NotNil(s.amqp.stopping)

	s.amqp.channel = s.cMock

	s.deliveries = make(chan amqp.Delivery)
}

func (s *TestAMQPSuite) SetupTest() {
	var err error
	select {
	case err = <-s.stop:
		break
	default:
	}
	s.Require().NoError(err)

	s.cMock.Calls = []mock.Call{}
	s.cMock.ExpectedCalls = []*mock.Call{}
	s.hMock.Calls = []mock.Call{}
	s.hMock.ExpectedCalls = []*mock.Call{}
}

func (s *TestAMQPSuite) Test00_StopBeforeStart() {
	s.Equal(NotStartedError, s.amqp.Stop())
}

func (s *TestAMQPSuite) Test01_ConsumeError() {
	expectedError := errors.New("testing")
	s.cMock.On("Consume", s.amqp.queue, mock.Anything, false, false, false, false, s.amqp.args).Return(s.deliveries, expectedError).Once()

	c := make(chan error)
	go func() { c <- s.amqp.Start() }() // Prevent being stuck.

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	var err error
	select {
	case err = <-c:
	case <-ticker.C:
	}
	s.Equal(expectedError, err)
}

func (s *TestAMQPSuite) Test02_Start() {
	s.cMock.On("Consume", s.amqp.queue, mock.Anything, false, false, false, false, s.amqp.args).Return(s.deliveries, nil).Once()

	s.amqp.ctx = nil
	s.amqp.cancel = nil

	s.stop = make(chan error)
	go func(c chan<- error, a *AMQP) { c <- a.Start() }(s.stop, s.amqp)

	s.amqp.waitStart()
}

func (s *TestAMQPSuite) Test03_StartAfterStart() {
	c := make(chan error)
	go func() { c <- s.amqp.Start() }() // Prevent being stuck.

	ticker := time.NewTicker(time.Second * 2)
	defer ticker.Stop()

	var err error
	select {
	case err = <-c:
	case <-ticker.C:
	}
	s.Equal(AlreadyStartedError, err)
}

func (s *TestAMQPSuite) Test04_HandleCreate() {
	expectedBody := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)

	expectedDelivery := amqp.Delivery{
		Acknowledger: s.cMock,
		Headers: amqp.Table{
			"X-Model": "scale_team",
			"X-Event": "create",
		},
		MessageCount: 0,
		DeliveryTag:  42,
		Body:         expectedBody,
	}

	s.hMock.On("HandleCreate", mock.Anything, expectedBody).Return(nil).Once()
	s.cMock.On("Ack", expectedDelivery.DeliveryTag, false).Return(nil).Once()

	go func() { s.deliveries <- expectedDelivery }()
	s.cMock.Wait()
}

func (s *TestAMQPSuite) Test05_HandleCreate_Reject() {
	expectedBody := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)

	expectedDelivery := amqp.Delivery{
		Acknowledger: s.cMock,
		Headers: amqp.Table{
			"X-Model": "scale_team",
			"X-Event": "create",
		},
		MessageCount: 0,
		DeliveryTag:  42,
		Body:         expectedBody,
	}

	s.hMock.On("HandleCreate", mock.Anything, expectedBody).Return(errors.New("testing")).Once()
	s.cMock.On("Reject", expectedDelivery.DeliveryTag, false).Return(nil).Once()

	go func() { s.deliveries <- expectedDelivery }()
	s.cMock.Wait()
}

func (s *TestAMQPSuite) Test06_HandleUpdate() {
	expectedBody := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)

	expectedDelivery := amqp.Delivery{
		Acknowledger: s.cMock,
		Headers: amqp.Table{
			"X-Model": "scale_team",
			"X-Event": "update",
		},
		MessageCount: 0,
		DeliveryTag:  42,
		Body:         expectedBody,
	}

	s.hMock.On("HandleUpdate", mock.Anything, expectedBody).Return(nil).Once()
	s.cMock.On("Ack", expectedDelivery.DeliveryTag, false).Return(nil).Once()

	go func() { s.deliveries <- expectedDelivery }()
	s.cMock.Wait()
}

func (s *TestAMQPSuite) Test07_HandleDestroy() {
	expectedBody := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)

	expectedDelivery := amqp.Delivery{
		Acknowledger: s.cMock,
		Headers: amqp.Table{
			"X-Model": "scale_team",
			"X-Event": "destroy",
		},
		MessageCount: 0,
		DeliveryTag:  42,
		Body:         expectedBody,
	}

	s.hMock.On("HandleDestroy", mock.Anything, expectedBody).Return(nil).Once()
	s.cMock.On("Ack", expectedDelivery.DeliveryTag, false).Return(nil).Once()

	go func() { s.deliveries <- expectedDelivery }()
	s.cMock.Wait()
}

func (s *TestAMQPSuite) Test08_BadModel() {
	expectedBody := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)

	expectedDelivery := amqp.Delivery{
		Acknowledger: s.cMock,
		Headers: amqp.Table{
			"X-Model": "techno",
			"X-Event": "destroy",
		},
		MessageCount: 0,
		DeliveryTag:  42,
		Body:         expectedBody,
	}

	s.cMock.On("Ack", expectedDelivery.DeliveryTag, false).Return(nil).Once()

	go func() { s.deliveries <- expectedDelivery }()
	s.cMock.Wait()
}

func (s *TestAMQPSuite) Test09_BadEvent() {
	expectedBody := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)

	expectedDelivery := amqp.Delivery{
		Acknowledger: s.cMock,
		Headers: amqp.Table{
			"X-Model": "scale_team",
			"X-Event": "casser",
		},
		MessageCount: 0,
		DeliveryTag:  42,
		Body:         expectedBody,
	}

	s.cMock.On("Ack", expectedDelivery.DeliveryTag, false).Return(nil).Once()

	go func() { s.deliveries <- expectedDelivery }()
	s.cMock.Wait()
}

func (s *TestAMQPSuite) Test10_BadEvent() {
	expectedBody := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)

	expectedDelivery := amqp.Delivery{
		Acknowledger: s.cMock,
		Headers: amqp.Table{
			"X-Model": "scale_team",
		},
		MessageCount: 0,
		DeliveryTag:  42,
		Body:         expectedBody,
	}

	s.cMock.On("Ack", expectedDelivery.DeliveryTag, false).Return(nil).Once()

	go func() { s.deliveries <- expectedDelivery }()
	s.cMock.Wait()
}

func (s *TestAMQPSuite) TearDownTest() {
	s.cMock.AssertExpectations(s.T())
	s.hMock.AssertExpectations(s.T())
}

func (s *TestAMQPSuite) TearDownSuite() {
	s.NoError(s.amqp.Stop())

	ticker := time.NewTicker(time.Second * 20)
	defer ticker.Stop()

	var err error
	select {
	case err = <-s.stop:
		break
	case <-ticker.C:
	}
	s.Equal(ConsumerStoppedError, err)
}
