package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gustavobelfort/42-jitsi/internal/config"
	"github.com/gustavobelfort/42-jitsi/internal/consumers"
	amqp2 "github.com/gustavobelfort/42-jitsi/internal/consumers/amqp"
	"github.com/gustavobelfort/42-jitsi/internal/db"
	"github.com/gustavobelfort/42-jitsi/internal/handler"
	"github.com/gustavobelfort/42-jitsi/internal/intra"
	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
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
	client, err := intra.NewClient(config.Conf.Intra.AppID, config.Conf.Intra.AppSecret, http.DefaultClient)
	if err != nil {
		logrus.Fatalf("could not initiate intra api client: %v", err)
	}

	conn, err := amqp.Dial(config.Conf.RabbitMQ.URL())
	if err != nil {
		logrus.Fatalf("could not connect to rabbitmq: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("could not initiate rabbitmq channel: %v", err)
	}

	hdl := handler.NewScaleTeamHandler(client, db.GlobalDB)
	consumer := amqp2.NewAMQP(channel, config.Conf.RabbitMQ.Queue, nil, hdl, config.Conf.Timeout)

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
			if err == amqp2.ConsumerStoppedError {
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
