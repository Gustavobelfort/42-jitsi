package tasks

import (
	"context"

	"github.com/gustavobelfort/42-jitsi/internal/config"
	"github.com/gustavobelfort/42-jitsi/internal/db"
	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/gustavobelfort/42-jitsi/internal/slack"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type tasksHandler struct {
	db *gorm.DB

	scaleTeamManager db.ScaleTeamManager
	userManager      db.UserManager

	client slack.SlackThat
}

type TasksHandler interface {
	Notify()
}

func NewTasksHandler(client slack.SlackThat, dbInstance *gorm.DB) TasksHandler {
	return &tasksHandler{
		db: dbInstance,

		scaleTeamManager: db.NewScaleTeamManager(dbInstance),
		userManager:      db.NewUserManager(dbInstance),
		client:           client,
	}
}

func (handler *tasksHandler) Notify() {
	logger := logging.ContextLog(context.Background(), logrus.StandardLogger())

	logger.Info("getting notifiable scale teams")
	scaleTeams, err := handler.getNotifiableScaleTeams()
	if err != nil {
		logger.Warnf("error getting notifiable scale teams: %v", err)
		return
	}

	if len(scaleTeams) == 0 {
		logger.Debugf("no scale teams to be notified")
		return
	}

	logger.Debugf("%d scale teams found to be notified", len(scaleTeams))
	for _, scaleTeam := range scaleTeams {
		scaleTeamID := scaleTeam.GetID()

		logins, err := handler.getScaleTeamUserLogins(scaleTeamID)
		if err != nil {
			logger.Warnf("error getting scale team logins: %v", err)
			return
		}

		err = handler.client.SendNotification(scaleTeamID, logins)
		if err != nil {
			logger.WithField("scale team", scaleTeamID).Warnf("error sending notification to the scale team: %v", err)
			return
		}

		scaleTeam.SetNotified(true)
		if err := scaleTeam.Save(handler.db); err != nil {
			logger.Warnf("error updating notified fild of the scale team: %v", err)
			return
		}
	}
	logger.Info("sucessfully notified scale teams")
}

func (handler *tasksHandler) getScaleTeamUserLogins(scaleTeamID int) ([]string, error) {
	var logins []string

	users, err := handler.userManager.Get(handler.db, db.UserScaleTeamOption(scaleTeamID))
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		login := user.GetLogin()
		logins = append(logins, login)
	}

	return logins, nil
}

func (handler *tasksHandler) getNotifiableScaleTeams() ([]db.ScaleTeam, error) {
	scaleTeams, err := handler.scaleTeamManager.Get(handler.db,
		db.ScaleTeamNotifiedOption(false),
		db.ScaleTeamBeginAtInOption(config.Conf.WarnBefore),
	)
	if err != nil {
		return nil, err
	}
	return scaleTeams, nil
}
