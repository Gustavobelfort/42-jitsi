package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

// UserStatus defines if a userModel is either corrector or part of the correcteds.
type UserStatus string

// UserStatus constant values.
var (
	Corrected UserStatus = "corrected"
	Corrector UserStatus = "corrector"
)

// ScaleTeamManager will be a wrapper to manage ScaleTeams in the database.
//
// It shall be used by a constant "GlobalScaleTeamManager".
type ScaleTeamManager interface {
	Create(tx *gorm.DB, id int, beginAt time.Time, notified bool) (ScaleTeam, error)
	Update(tx *gorm.DB, scaleTeam ScaleTeam) error
	Delete(tx *gorm.DB, scaleTeam ScaleTeam) error
	Get(tx *gorm.DB, options ...GetOption) ([]ScaleTeam, error)

	DB() *gorm.DB
}

// UserManager will be a wrapper to manage Users in the database.
//
// It shall be used by a constant "GlobalUserManager".
type UserManager interface {
	Create(tx *gorm.DB, scaleTeamID int, login string, status UserStatus) (User, error)
	Update(tx *gorm.DB, user User) error
	Delete(tx *gorm.DB, user User) error
	Get(tx *gorm.DB, options ...GetOption) ([]User, error)

	DB() *gorm.DB
}

// ManagedModel is a base interface for managed data models.
type ManagedModel interface {
	// Delete the data inheriting this model.
	Delete(tx *gorm.DB) error

	// Save the data inheriting this model.
	Save(tx *gorm.DB) error
}

// ScaleTeam and manages wraps the scale_teams records.
type ScaleTeam interface {
	GetID() int
	GetBeginAt() time.Time
	GetNotified() bool

	Get(tx *gorm.DB, options ...GetOption) ([]User, error)

	SetID(int)
	SetBeginAt(time.Time)
	SetNotified(bool)

	ManagedModel
}

// User wraps and manages the users records.
type User interface {
	GetID() int
	GetScaleTeamID() int
	GetLogin() string
	GetStatus() UserStatus

	SetScaleTeamID(int)
	SetLogin(string)
	SetStatus(UserStatus)

	ManagedModel
}
