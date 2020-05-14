package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gustavobelfort/42-jitsi/internal/config"
	"github.com/gustavobelfort/42-jitsi/internal/consumers"
	"github.com/gustavobelfort/42-jitsi/internal/db"
	"github.com/gustavobelfort/42-jitsi/internal/intra"
	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/gustavobelfort/42-jitsi/internal/scheduler"
	"github.com/gustavobelfort/42-jitsi/internal/slack"
	"github.com/gustavobelfort/42-jitsi/internal/tasks"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := config.Initiate(); err != nil {
		logrus.Fatalf("could not load configuration: %v", err)
	}
	logging.Initiate()
	if err := db.Init(); err != nil {
		logrus.Fatalf("could not connect to the db: %v", err)
	}
}

func main() {

	iClient, err := intra.NewClient(config.Conf.Intra.AppID, config.Conf.Intra.AppSecret, http.DefaultClient)
	if err != nil {
		logrus.Fatalf("could not initiate intra api client: %v", err)
	}

	sClient, err := slack.New(iClient, config.Conf.SlackThat.URL)
	if err != nil {
		logrus.Fatalf("could not initiate slack_that client: %v", err)
	}

	tHdl := tasks.NewTasksHandler(sClient, db.GlobalDB)
	tasks := scheduler.Task{
		Task:     tHdl.Notify,
		Interval: config.Conf.WarnBefore,
	}

	scheduler, err := scheduler.New([]scheduler.Task{tasks})
	if err != nil {
		logrus.Fatalf("could not create scheduler: %v", err)
	}

	waitForShutdown(scheduler)
}

func waitForShutdown(consumer consumers.Consumer) {
	interruptChan := make(chan os.Signal)
	isDown := make(chan struct{})
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer func() { isDown <- struct{}{} }()
		err := consumer.Start()
		if err != nil {
			logrus.Error(err)
			return
		}
	}()
	<-interruptChan
	logging.LogError(logrus.StandardLogger(), consumer.Stop(), "while shutting down")
	<-isDown
}
