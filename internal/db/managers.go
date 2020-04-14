package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

/*
 * Scale Teams Manager
 */

type scaleTeamManager struct {
	db *gorm.DB
}

func (stManager *scaleTeamManager) CreateWithTransaction(tx *gorm.DB, id int, beginAt time.Time, notified bool) (ScaleTeam, error) {
	scaleTeam := &scaleTeamModel{
		ID:       id,
		BeginAt:  beginAt,
		Notified: notified,

		scaleTeamManager: stManager,
		userManager:      &userManager{db: stManager.db},
	}

	if err := tx.Create(scaleTeam).Error; err != nil {
		return nil, err
	}
	return scaleTeam, nil
}

func (stManager *scaleTeamManager) Create(id int, beginAt time.Time, notified bool) (ScaleTeam, error) {
	tx := stManager.db.Begin()
	defer tx.RollbackUnlessCommitted()
	scaleTeam, err := stManager.CreateWithTransaction(tx, id, beginAt, notified)
	if err != nil {
		return nil, err
	}
	return scaleTeam, tx.Commit().Error
}

func (stManager *scaleTeamManager) UpdateWithTransaction(tx *gorm.DB, scaleTeam ScaleTeam) error {
	return tx.Save(scaleTeam).Error
}

func (stManager *scaleTeamManager) Update(scaleTeam ScaleTeam) error {
	tx := stManager.db.Begin()
	defer tx.RollbackUnlessCommitted()
	if err := stManager.UpdateWithTransaction(tx, scaleTeam); err != nil {
		return err
	}
	return tx.Commit().Error
}

func (stManager *scaleTeamManager) DeleteWithTransaction(tx *gorm.DB, scaleTeam ScaleTeam) error {
	return tx.Delete(scaleTeam).Error
}

func (stManager *scaleTeamManager) Delete(scaleTeam ScaleTeam) error {
	tx := stManager.db.Begin()
	defer tx.RollbackUnlessCommitted()
	if err := stManager.DeleteWithTransaction(tx, scaleTeam); err != nil {
		return err
	}
	return stManager.db.Delete(scaleTeam).Error
}

func (stManager *scaleTeamManager) GetWithTransaction(tx *gorm.DB, options ...GetOption) ([]ScaleTeam, error) {
	for _, opt := range options {
		tx = opt(tx)
	}
	var scaleTeams []scaleTeamModel

	if err := tx.Find(&scaleTeams).Error; err != nil {
		return nil, err
	}

	returned := make([]ScaleTeam, len(scaleTeams))
	for i, scaleTeam := range scaleTeams {
		scaleTeam.scaleTeamManager = stManager
		scaleTeam.userManager = &userManager{db: stManager.db}
		returned[i] = &scaleTeam
	}
	return returned, nil

}

func (stManager *scaleTeamManager) Get(options ...GetOption) ([]ScaleTeam, error) {
	return stManager.GetWithTransaction(stManager.db, options...)
}

/*
 * Users Manager
 */

type userManager struct {
	db *gorm.DB
}

func (uManager *userManager) CreateWithTransaction(tx *gorm.DB, scaleTeamID int, login string, status UserStatus) (User, error) {
	user := &userModel{
		ScaleTeamID: scaleTeamID,
		Login:       login,
		Status:      status,

		scaleTeamManager: &scaleTeamManager{db: uManager.db},
		userManager:      uManager,
	}
	return user, tx.Create(user).Error
}

func (uManager *userManager) Create(scaleTeamID int, login string, status UserStatus) (User, error) {
	tx := uManager.db.Begin()
	defer tx.RollbackUnlessCommitted()
	user, err := uManager.CreateWithTransaction(tx, scaleTeamID, login, status)
	if err != nil {
		return nil, err
	}
	return user, tx.Commit().Error
}

func (uManager *userManager) UpdateWithTransaction(tx *gorm.DB, user User) error {
	return tx.Save(user).Error
}

func (uManager *userManager) Update(user User) error {
	tx := uManager.db.Begin()
	defer tx.RollbackUnlessCommitted()
	if err := uManager.UpdateWithTransaction(tx, user); err != nil {
		return err
	}
	return tx.Commit().Error
}

func (uManager *userManager) DeleteWithTransaction(tx *gorm.DB, user User) error {
	return tx.Delete(user).Error
}

func (uManager *userManager) Delete(user User) error {
	tx := uManager.db.Begin()
	defer tx.RollbackUnlessCommitted()
	if err := uManager.DeleteWithTransaction(tx, user); err != nil {
		return err
	}
	return tx.Commit().Error
}

func (uManager *userManager) GetWithTransaction(tx *gorm.DB, options ...GetOption) ([]User, error) {
	for _, opt := range options {
		tx = opt(tx)
	}
	var users []userModel

	if err := tx.Find(&users).Error; err != nil {
		return nil, err
	}

	returned := make([]User, len(users))
	for i, user := range users {
		user.scaleTeamManager = &scaleTeamManager{db: uManager.db}
		user.userManager = uManager
		returned[i] = &users[i]
	}
	return returned, nil
}

func (uManager *userManager) Get(options ...GetOption) ([]User, error) {
	return uManager.GetWithTransaction(uManager.db, options...)
}
