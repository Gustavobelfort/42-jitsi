package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

type scaleTeamModel struct {
	ID       int `gorm:"primary_key;auto_increment:false"`
	BeginAt  time.Time
	Notified bool        `gorm:"default:false"`
	Users    []userModel `gorm:"foreignkey:ScaleTeamID"`

	userManager      UserManager      `gorm:"-"`
	scaleTeamManager ScaleTeamManager `gorm:"-"`
}

func (scaleTeamModel) TableName() string {
	return "scale_teams"
}

func (scaleTeam *scaleTeamModel) GetID() int {
	return scaleTeam.ID
}

func (scaleTeam *scaleTeamModel) GetBeginAt() time.Time {
	return scaleTeam.BeginAt
}

func (scaleTeam *scaleTeamModel) GetNotified() bool {
	return scaleTeam.Notified
}

func (scaleTeam *scaleTeamModel) Get(tx *gorm.DB, options ...GetOption) ([]User, error) {
	options = append(options, UserScaleTeamOption(scaleTeam.ID))
	return GlobalUserManager.Get(tx, options...)
}

func (scaleTeam *scaleTeamModel) SetID(id int) {
	scaleTeam.ID = id
}

func (scaleTeam *scaleTeamModel) SetBeginAt(beginAt time.Time) {
	scaleTeam.BeginAt = beginAt
}

func (scaleTeam *scaleTeamModel) SetNotified(notified bool) {
	scaleTeam.Notified = notified
}

func (scaleTeam *scaleTeamModel) Save(tx *gorm.DB) error {
	return scaleTeam.scaleTeamManager.Update(tx, scaleTeam)
}

func (scaleTeam *scaleTeamModel) Delete(tx *gorm.DB) error {
	return scaleTeam.scaleTeamManager.Delete(tx, scaleTeam)
}

type userModel struct {
	ID          int `gorm:"primary_key"`
	ScaleTeamID int
	Login       string     `gorm:"varchar(32)"`
	Status      UserStatus `gorm:"varchar(32)"`

	userManager      UserManager      `gorm:"-"`
	scaleTeamManager ScaleTeamManager `gorm:"-"`
}

func (userModel) TableName() string {
	return "users"
}

func (user *userModel) GetID() int {
	return user.ID
}

func (user *userModel) GetScaleTeamID() int {
	return user.ScaleTeamID
}

func (user *userModel) GetLogin() string {
	return user.Login
}

func (user *userModel) GetStatus() UserStatus {
	return user.Status
}

func (user *userModel) SetScaleTeamID(scaleTeamID int) {
	user.ScaleTeamID = scaleTeamID
}

func (user *userModel) SetLogin(login string) {
	user.Login = login
}

func (user *userModel) SetStatus(status UserStatus) {
	user.Status = status
}

func (user *userModel) Save(tx *gorm.DB) error {
	return user.userManager.Update(tx, user)
}

func (user *userModel) Delete(tx *gorm.DB) error {
	return user.userManager.Delete(tx, user)
}
