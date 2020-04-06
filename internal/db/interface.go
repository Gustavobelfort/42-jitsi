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

// ScaleTeamManagerInterface will be a wrapper to manage ScaleTeams in the database.
//
// It shall be used by a constant "ScaleTeamManager".
type ScaleTeamManagerInterface interface {
	Create(id int, beginAt time.Time, notified bool) (ScaleTeam, error)
	CreateWithTransaction(tx *gorm.DB, id int, beginAt time.Time, notified bool) (ScaleTeam, error)

	Update(scaleTeam ScaleTeam) error
	UpdateWithTransaction(tx *gorm.DB, scaleTeam ScaleTeam) error

	Delete(scaleTeam ScaleTeam) error
	DeleteWithTransaction(tx *gorm.DB, scaleTeam ScaleTeam) error

	Get(options ...GetOption) ([]ScaleTeam, error)
	GetWithTransaction(tx *gorm.DB, options ...GetOption) ([]ScaleTeam, error)
}

// UserManagerInterface will be a wrapper to manage Users in the database.
//
// It shall be used by a constant "UserManager".
type UserManagerInterface interface {
	Create(scaleTeamID int, login string, status UserStatus) (User, error)
	CreateWithTransaction(tx *gorm.DB, scaleTeamID int, login string, status UserStatus) (User, error)

	Update(scaleTeam User) error
	UpdateWithTransaction(tx *gorm.DB, user User) error

	Delete(user User) error
	DeleteWithTransaction(tx *gorm.DB, user User) error

	Get(options ...GetOption) ([]User, error)
	GetWithTransaction(tx *gorm.DB, options ...GetOption) ([]User, error)
}

// ManagedModel is a base interface for managed data models.
type ManagedModel interface {
	// Delete the data inheriting this model.
	Delete() error
	DeleteWithTransaction(tx *gorm.DB) error

	// Save the data inheriting this model.
	Save() error
	SaveWithTransaction(tx *gorm.DB) error
}

// ScaleTeam and manages wraps the scale_teams records.
type ScaleTeam interface {
	GetID() int
	GetBeginAt() time.Time
	GetNotified() bool

	GetUsers(...GetOption) ([]User, error)

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
