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

// NewScaleTeamManager returns a new manager with the passed GlobalDB object.
func NewScaleTeamManager(db *gorm.DB) ScaleTeamManager {
	return &scaleTeamManager{db: db}
}

// Returns the underlying database object.
func (stManager *scaleTeamManager) DB() *gorm.DB {
	return stManager.db
}

func (stManager *scaleTeamManager) Create(tx *gorm.DB, id int, beginAt time.Time, notified bool) (ScaleTeam, error) {
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

func (stManager *scaleTeamManager) Update(tx *gorm.DB, scaleTeam ScaleTeam) error {
	return tx.Save(scaleTeam).Error
}

func (stManager *scaleTeamManager) Delete(tx *gorm.DB, scaleTeam ScaleTeam) error {
	return tx.Delete(scaleTeam).Error
}

func (stManager *scaleTeamManager) Get(tx *gorm.DB, options ...GetOption) ([]ScaleTeam, error) {
	for _, opt := range options {
		tx = opt(tx)
	}
	var scaleTeams []scaleTeamModel

	if err := tx.Find(&scaleTeams).Error; err != nil {
		return nil, err
	}

	returned := make([]ScaleTeam, len(scaleTeams))
	for i := range scaleTeams {
		scaleTeams[i].scaleTeamManager = stManager
		scaleTeams[i].userManager = &userManager{db: stManager.db}
		returned[i] = &scaleTeams[i]
	}
	return returned, nil

}

/*
 * Users Manager
 */

type userManager struct {
	db *gorm.DB
}

// NewUserManager returns a new manager with the passed GlobalDB object.
func NewUserManager(db *gorm.DB) UserManager {
	return &userManager{db: db}
}

// Returns the underlying database object.
func (uManager *userManager) DB() *gorm.DB {
	return uManager.db
}

func (uManager *userManager) Create(tx *gorm.DB, scaleTeamID int, login string, status UserStatus) (User, error) {
	user := &userModel{
		ScaleTeamID: scaleTeamID,
		Login:       login,
		Status:      status,

		scaleTeamManager: &scaleTeamManager{db: uManager.db},
		userManager:      uManager,
	}
	if err := tx.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (uManager *userManager) Update(tx *gorm.DB, user User) error {
	return tx.Save(user).Error
}

func (uManager *userManager) Delete(tx *gorm.DB, user User) error {
	return tx.Delete(user).Error
}

func (uManager *userManager) Get(tx *gorm.DB, options ...GetOption) ([]User, error) {
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
