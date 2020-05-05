package router

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gustavobelfort/42-jitsi/internal/consumers"
	"github.com/gustavobelfort/42-jitsi/internal/handler"
	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/sirupsen/logrus"
)

// Router will consume the scale teams from the intranet's webhooks by exposing a endpoint.
type Router struct {
	engine *gin.Engine
	server *http.Server

	handler handler.ScaleTeamHandler

	registries map[string]string

	mu     *sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
}

// NewRouter returns a new router consumer.
func NewRouter(server *http.Server, hdl handler.ScaleTeamHandler, registries map[string]string, prefix string) consumers.Consumer {
	router := &Router{
		engine: nil,
		server: server,

		handler: hdl,

		registries: registries,

		mu: new(sync.Mutex),
	}
	router.setupEngine(prefix)
	return router
}

func (r *Router) setContext(listener net.Listener) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.ctx != nil {
		return AlreadyStartedError
	}

	r.ctx, r.cancel = context.WithCancel(context.Background())
	r.ctx = logging.ContextWithFields(r.ctx, logrus.Fields{"consumer": "gin", "listening": listener.Addr()})
	r.server.RegisterOnShutdown(r.cancel)
	return nil
}

func (r *Router) start(listener net.Listener) error {
	if err := r.setContext(listener); err != nil {
		return err
	}
	logging.ContextLog(r.ctx, logrus.StandardLogger()).Info("start listening")
	r.server.Handler = r.engine

	return r.server.Serve(listener)
}

// Start tries to start the consumer. If it's already started, it returns `AlreadyStartedError`.
func (r *Router) Start() error {
	addr := r.server.Addr
	if addr == "" {
		addr = ":http"
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	return r.start(listener)
}

// Stop faithfully shuts down the consumer with a timeout of 20 seconds.
func (r *Router) Stop() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	logging.ContextLog(r.ctx, logrus.StandardLogger()).WithField("timeout", time.Second*20).Info("shutting down consumer")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()
	return r.server.Shutdown(ctx)
}
