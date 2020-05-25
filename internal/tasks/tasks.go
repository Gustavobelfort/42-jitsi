package tasks

import (
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
	logrus.Debug("getting notifiable scale teams")
	scaleTeams, err := handler.getNotifiableScaleTeams()
	if err != nil {
		logrus.WithError(err).Errorf("error getting notifiable scale teams: %v", err)
		return
	}

	if len(scaleTeams) == 0 {
		logrus.Debugf("no scale teams to be notified")
		return
	}

	logrus.Infof("found %d scale teams to notify", len(scaleTeams))
	for _, scaleTeam := range scaleTeams {
		scaleTeamID := scaleTeam.GetID()
		ctxlogger := logrus.WithField("scale_team_id", scaleTeamID)

		logins, err := handler.getScaleTeamUserLogins(scaleTeamID)
		if err != nil {
			logging.LogError(ctxlogger, err, "getting scale team users' logins")
			continue
		}

		err = handler.client.SendNotification(scaleTeamID, logins)
		if err != nil {
			logging.LogError(ctxlogger, err, "sending notification to the scale team")
			continue
		}

		scaleTeam.SetNotified(true)
		if err := scaleTeam.Save(handler.db); err != nil {
			logging.LogError(ctxlogger, err, "updating scale team notified field")
			continue
		}
		ctxlogger.Info("successfully notified scale team")
	}
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
