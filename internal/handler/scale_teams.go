package handler

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/gustavobelfort/42-jitsi/internal/db"
	"github.com/gustavobelfort/42-jitsi/internal/intra"
	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/gustavobelfort/42-jitsi/internal/utils"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type scaleTeamHandler struct {
	db *gorm.DB

	scaleTeamManager db.ScaleTeamManager
	userManager      db.UserManager

	client intra.Client
}

// NewScaleTeamHandler returns a new handler that will handle scale teams payloads with the given client and db managers.
func NewScaleTeamHandler(client intra.Client, dbInstance *gorm.DB) ScaleTeamHandler {
	return &scaleTeamHandler{
		db: dbInstance,

		scaleTeamManager: db.NewScaleTeamManager(dbInstance),
		userManager:      db.NewUserManager(dbInstance),
		client:           client,
	}
}

func (handler *scaleTeamHandler) interpretData(ctx context.Context, data []byte, logger *logrus.Entry) (*scaleTeam, error) {
	st := &scaleTeam{}

	logger.Info("parsing webhook's payload")
	err := utils.WrapContext(ctx, func() error {
		return json.Unmarshal(data, &st)
	})
	if err != nil {
		return nil, err
	}

	logger.Info("getting scale team's corrected members")
	if st.Correcteds, err = handler.client.GetTeamMembers(ctx, st.TeamID); err != nil {
		return nil, err
	}

	return st, nil
}

func (handler *scaleTeamHandler) insertInDB(tx *gorm.DB, st *scaleTeam, logger *logrus.Entry) error {
	defer tx.RollbackUnlessCommitted()

	logger.Info("creating scale team's record")
	stRecord, err := handler.scaleTeamManager.Create(tx, st.ID, st.BeginAt, false)
	if err != nil {
		return err
	}
	logger.WithField("login", st.Corrector).Info("creating scale team's corrector record")
	if _, err = handler.userManager.Create(tx, stRecord.GetID(), st.Corrector, db.Corrector); err != nil {
		return err
	}
	for _, login := range st.Correcteds {
		logger.WithField("login", login).Info("creating scale team's corrected record")
		if _, err = handler.userManager.Create(tx, stRecord.GetID(), login, db.Corrected); err != nil {
			return err
		}
	}
	return tx.Commit().Error
}

func (handler *scaleTeamHandler) HandleCreate(ctx context.Context, data []byte) error {
	logger := logging.ContextLog(ctx, logrus.StandardLogger())

	st, err := handler.interpretData(ctx, data, logger)
	if err != nil {
		return err
	}

	return handler.insertInDB(handler.db.BeginTx(ctx, &sql.TxOptions{}), st, logger.WithField("scale_team_id", st.ID))
}

func (handler *scaleTeamHandler) updateInDB(tx *gorm.DB, st *scaleTeam, logger *logrus.Entry) error {
	defer tx.RollbackUnlessCommitted()

	logger.Info("getting corresponding scale team's record")
	stRecords, err := handler.scaleTeamManager.Get(tx, db.ScaleTeamIDOption(st.ID))
	if err != nil {
		return err
	}

	if len(stRecords) == 0 {
		logger.WithField("error", NotInDBError).Warnf("trying to update scale team: %v", NotInDBError)
		return handler.insertInDB(tx, st, logger)
	}

	if st.BeginAt.Equal(stRecords[0].GetBeginAt()) {
		logger.Info("scale team's begin_at did not change")
		return nil
	}

	logger.Info("updating scale team's record")
	logger.Debugf("setting begin_at to: %v", st.BeginAt)
	stRecords[0].SetBeginAt(st.BeginAt)
	logger.Debugf("setting notified to: %v", false)
	stRecords[0].SetNotified(false)
	if err := stRecords[0].Save(tx); err != nil {
		return err
	}
	return tx.Commit().Error
}

func (handler *scaleTeamHandler) HandleUpdate(ctx context.Context, data []byte) error {
	logger := logging.ContextLog(ctx, logrus.StandardLogger())

	st, err := handler.interpretData(ctx, data, logger)
	if err != nil {
		return err
	}

	return handler.updateInDB(handler.db.BeginTx(ctx, &sql.TxOptions{}), st, logger)
}

func (handler *scaleTeamHandler) deleteFromDB(tx *gorm.DB, id int, logger *logrus.Entry) error {
	defer tx.RollbackUnlessCommitted()

	logger.Info("getting corresponding scale team's records")
	stRecords, err := handler.scaleTeamManager.Get(tx, db.ScaleTeamIDOption(id))
	if err != nil {
		return err
	}

	if len(stRecords) == 0 {
		return logging.WithLog(NotInDBError, logrus.WarnLevel, logrus.Fields{"scale_team_id": id})
	}

	logger.Info("deleting scale team's records")
	for _, record := range stRecords {
		if err := record.Delete(tx); err != nil {
			return err
		}
	}

	return tx.Commit().Error
}

func (handler *scaleTeamHandler) HandleDestroy(ctx context.Context, data []byte) error {
	logger := logging.ContextLog(ctx, logrus.StandardLogger())
	st := make(map[string]interface{})
	logger.Info("parsing webhook's payload")
	err := utils.WrapContext(ctx, func() error {
		if err := json.Unmarshal(data, &st); err != nil {
			return err
		}
		logger.Debug("checking if id is present")
		if _, ok := st["id"]; !ok {
			return &MissingFieldsError{missing: []string{"id"}}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return handler.deleteFromDB(handler.db.BeginTx(ctx, &sql.TxOptions{}), int(st["id"].(float64)), logger)
}
