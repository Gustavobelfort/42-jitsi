package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	if err := config.Initiate(); err != nil {
		logrus.Fatalf("could not load configuration: %v", err)
	}
	if err := db.Init(); err != nil {
		logrus.Fatalf("could not connect to the db: %v", err)
	}
}

func main() {
	server := &http.Server{
		Addr:         "0.0.0.0:5000",
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 10,
	}

	client, err := intra.NewClient(config.Conf.Intra.AppID, config.Conf.Intra.AppSecret, http.DefaultClient)
	if err != nil {
		logrus.Fatalf("could not initiate intra api client: %v", err)
	}

	hdl := handler.NewScaleTeamHandler(client, db.GlobalDB)
	consumer := router.NewRouter(server, hdl, config.Conf.Intra.Webhooks, "/")

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
			logrus.Error(err)
		}
	}()
	<-interruptChan
	logging.LogError(logrus.StandardLogger(), consumer.Stop(), "while shutting down")
	<-isDown
}
