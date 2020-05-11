package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gustavobelfort/42-jitsi/internal/config"
	"github.com/gustavobelfort/42-jitsi/internal/consumers"
	"github.com/gustavobelfort/42-jitsi/internal/consumers/router"
	"github.com/gustavobelfort/42-jitsi/internal/db"
	"github.com/gustavobelfort/42-jitsi/internal/handler"
	"github.com/gustavobelfort/42-jitsi/internal/intra"
	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/sirupsen/logrus"
)

func init() {
	config.AddRequired("intra.app_id", "intra.app_secret", "intra.webhooks")
	if err := config.Initiate(); err != nil {
		logrus.WithError(err).Fatalf("could not load configuration: %v", err)
	}
	logging.Initiate()
	if err := db.Init(); err != nil {
		logrus.WithError(err).Fatalf("could not connect to the db: %v", err)
	}
}

func main() {
	if hub := sentry.CurrentHub(); hub.Client() != nil {
		defer hub.Flush(time.Second * 5)
	}
	server := &http.Server{
		Addr:         config.Conf.HTTPAddr,
		ReadTimeout:  config.Conf.Timeout * 2,
		WriteTimeout: config.Conf.Timeout * 2,
		IdleTimeout:  config.Conf.Timeout * 2,
	}

	iClient, err := intra.NewClient(config.Conf.Intra.AppID, config.Conf.Intra.AppSecret, http.DefaultClient)
	if err != nil {
		logrus.WithError(err).Fatalf("could not initiate intra api client: %v", err)
	}

	stHdl := handler.NewScaleTeamHandler(iClient, db.GlobalDB)
	consumer := router.NewRouter(server, stHdl, config.Conf.Intra.Webhooks, "/", config.Conf.Timeout)

	waitForShutdown(consumer)
}

func waitForShutdown(consumer consumers.Consumer) {
	interruptChan := make(chan os.Signal)
	isDown := make(chan struct{})
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer func() { isDown <- struct{}{} }()
		err := consumer.Start()
		if err != nil {
			if err == http.ErrServerClosed {
				logrus.Info(err)
				return
			}
			logrus.WithError(err).Error(err)
		}
	}()
	<-interruptChan
	logging.LogError(logrus.StandardLogger(), consumer.Stop(), "while shutting down")
	<-isDown
}
