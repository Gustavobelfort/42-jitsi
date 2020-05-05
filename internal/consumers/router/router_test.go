package router

import (
	"bytes"
	"context"
	"errors"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"golang.org/x/net/nettest"
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

func TestRouter(t *testing.T) {
	t.Run("TestRouter_StartStop", func(t *testing.T) {
		server := &http.Server{
			Addr: "127.0.0.1:8442",
		}
		hdl := &HandlerMock{}

		router := NewRouter(server, hdl, nil, "/")

		ch := make(chan error)
		go func() { ch <- router.Start() }()

		listener, err := nettest.NewLocalListener("tcp")
		if assert.NoError(t, err) {
			assert.Equal(t, AlreadyStartedError, router.(*Router).start(listener))
		}

		assert.Error(t, router.Start())

		assert.NoError(t, router.Stop())

		ticker := time.NewTicker(time.Second * 20)
		defer ticker.Stop()

		select {
		case err = <-ch:
			break
		case <-ticker.C:
		}
		assert.Equal(t, http.ErrServerClosed, err)
	})

	suite.Run(t, new(TestRouterSuite))
}

type TestRouterSuite struct {
	suite.Suite

	mock       *HandlerMock
	registries map[string]string

	listener net.Listener
	router   *Router

	stop chan error
}

func (s *TestRouterSuite) SetupSuite() {
	var err error
	s.listener, err = nettest.NewLocalListener("tcp")
	s.Require().NoError(err)

	s.mock = &HandlerMock{}
	s.registries = map[string]string{
		"scale_team.create":  "create_secret",
		"scale_team.update":  "update_secret",
		"scale_team.destroy": "destroy_secret",
		"scale_team.unknown": "unknown_secret",
		"unknown.unknown":    "unknown_secret",
	}

	s.router = NewRouter(&http.Server{
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 10,
	}, s.mock, s.registries, "/").(*Router)

	s.stop = make(chan error)
	go func(c chan<- error, r *Router) { c <- r.start(s.listener) }(s.stop, s.router)
}

func (s *TestRouterSuite) SetupTest() {
	var err error
	select {
	case err = <-s.stop:
		break
	default:
	}
	s.Require().NoError(err)

	s.mock.Calls = []mock.Call{}
	s.mock.ExpectedCalls = []*mock.Call{}
}

func (s *TestRouterSuite) Test00_CreateWebhook() {
	body := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)
	buffer := bytes.NewBuffer(body)

	request, err := http.NewRequest(http.MethodPost, "http://"+s.listener.Addr().String()+"/webhooks", buffer)
	s.Require().NoError(err)

	request.Header.Set("X-Model", "scale_team")
	request.Header.Set("X-Event", "create")
	request.Header.Set("X-Secret", s.registries["scale_team.create"])

	s.mock.On("HandleCreate", mock.Anything, body).Return(nil).Once()

	resp, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)

	s.Equal(http.StatusNoContent, resp.StatusCode)
}

func (s *TestRouterSuite) Test01_CreateWebhook_BadPayload() {
	body := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin",
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)
	buffer := bytes.NewBuffer(body)

	request, err := http.NewRequest(http.MethodPost, "http://"+s.listener.Addr().String()+"/webhooks", buffer)
	s.Require().NoError(err)

	request.Header.Set("X-Model", "scale_team")
	request.Header.Set("X-Event", "create")
	request.Header.Set("X-Secret", s.registries["scale_team.create"])

	expectedErr := logging.WithLog(errors.New("testing"), logrus.WarnLevel, nil)

	s.mock.On("HandleCreate", mock.Anything, body).Return(expectedErr).Once()

	resp, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)

	s.Equal(http.StatusBadRequest, resp.StatusCode)
}

func (s *TestRouterSuite) Test02_CreateWebhook_InternalServerError() {
	body := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)
	buffer := bytes.NewBuffer(body)

	request, err := http.NewRequest(http.MethodPost, "http://"+s.listener.Addr().String()+"/webhooks", buffer)
	s.Require().NoError(err)

	request.Header.Set("X-Model", "scale_team")
	request.Header.Set("X-Event", "create")
	request.Header.Set("X-Secret", s.registries["scale_team.create"])

	s.mock.On("HandleCreate", mock.Anything, body).Return(errors.New("testing")).Once()

	resp, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)

	s.Equal(http.StatusInternalServerError, resp.StatusCode)
}

func (s *TestRouterSuite) Test03_UpdateWebhook() {
	body := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)
	buffer := bytes.NewBuffer(body)

	request, err := http.NewRequest(http.MethodPost, "http://"+s.listener.Addr().String()+"/webhooks", buffer)
	s.Require().NoError(err)

	request.Header.Set("X-Model", "scale_team")
	request.Header.Set("X-Event", "update")
	request.Header.Set("X-Secret", s.registries["scale_team.update"])

	s.mock.On("HandleUpdate", mock.Anything, body).Return(nil).Once()

	resp, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)

	s.Equal(http.StatusNoContent, resp.StatusCode)
}

func (s *TestRouterSuite) Test04_UpdateWebhook_Unauthorized() {
	body := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)
	buffer := bytes.NewBuffer(body)

	request, err := http.NewRequest(http.MethodPost, "http://"+s.listener.Addr().String()+"/webhooks", buffer)
	s.Require().NoError(err)

	request.Header.Set("X-Model", "scale_team")
	request.Header.Set("X-Event", "update")
	request.Header.Set("X-Secret", "bad_secret")

	resp, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)

	s.Equal(http.StatusUnauthorized, resp.StatusCode)
}

func (s *TestRouterSuite) Test05_DestroyWebhook() {
	body := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)
	buffer := bytes.NewBuffer(body)

	request, err := http.NewRequest(http.MethodPost, "http://"+s.listener.Addr().String()+"/webhooks", buffer)
	s.Require().NoError(err)

	request.Header.Set("X-Model", "scale_team")
	request.Header.Set("X-Event", "destroy")
	request.Header.Set("X-Secret", s.registries["scale_team.destroy"])

	s.mock.On("HandleDestroy", mock.Anything, body).Return(nil).Once()

	resp, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)

	s.Equal(http.StatusNoContent, resp.StatusCode)
}

func (s *TestRouterSuite) Test06_UnknownModelWebhook() {
	body := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)
	buffer := bytes.NewBuffer(body)

	request, err := http.NewRequest(http.MethodPost, "http://"+s.listener.Addr().String()+"/webhooks", buffer)
	s.Require().NoError(err)

	request.Header.Set("X-Model", "unknown")
	request.Header.Set("X-Event", "unknown")
	request.Header.Set("X-Secret", s.registries["unknown.unknown"])

	resp, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)

	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *TestRouterSuite) Test07_UnknownEventWebhook() {
	body := []byte(`{
	"id": 21,
	"team": {"id": "42"},
	"user": {"login": "xlogin"},
	"begin_at": "2020-05-05T16:00:00.051Z"
}`)
	buffer := bytes.NewBuffer(body)

	request, err := http.NewRequest(http.MethodPost, "http://"+s.listener.Addr().String()+"/webhooks", buffer)
	s.Require().NoError(err)

	request.Header.Set("X-Model", "scale_team")
	request.Header.Set("X-Event", "unknown")
	request.Header.Set("X-Secret", s.registries["scale_team.unknown"])

	resp, err := http.DefaultClient.Do(request)
	s.Require().NoError(err)

	s.Equal(http.StatusNotFound, resp.StatusCode)
}

func (s *TestRouterSuite) TearDownTest() {
	s.mock.AssertExpectations(s.T())
}

func (s *TestRouterSuite) TearDownSuite() {
	s.NoError(s.router.Stop())

	ticker := time.NewTicker(time.Second * 20)
	defer ticker.Stop()

	var err error
	select {
	case err = <-s.stop:
		break
	case <-ticker.C:
	}
	s.Equal(http.ErrServerClosed, err)
}
