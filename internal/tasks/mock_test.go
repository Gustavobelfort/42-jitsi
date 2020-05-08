package tasks

import (
	"context"
	"database/sql"
	"time"

	"github.com/gustavobelfort/42-jitsi/internal/db"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"
)

type TxManagerMock struct {
	mock.Mock
}

func (m *TxManagerMock) BeginTx(ctx context.Context, options *sql.TxOptions) *gorm.DB {
	return m.Called(ctx, options).Get(0).(*gorm.DB)
}

type ClientMock struct {
	mock.Mock
}

func (m *ClientMock) GetHealth() (map[string]interface{}, error) {
	return nil, nil
}
func (m *ClientMock) SendNotification(scaleTeamID int, logins []string) error {
	return nil
}

type ScaleTeamManagerMock struct {
	mock.Mock
}

func (m *ScaleTeamManagerMock) DB() *gorm.DB {
	return m.Called().Get(0).(*gorm.DB)
}

func (m *ScaleTeamManagerMock) Create(tx *gorm.DB, id int, beginAt time.Time, notified bool) (db.ScaleTeam, error) {
	toReturn := m.Called(tx, id, beginAt, notified)
	return toReturn.Get(0).(db.ScaleTeam), toReturn.Error(1)
}

func (m *ScaleTeamManagerMock) Update(tx *gorm.DB, scaleTeam db.ScaleTeam) error {
	return m.Called(tx, scaleTeam).Error(0)
}

func (m *ScaleTeamManagerMock) Delete(tx *gorm.DB, scaleTeam db.ScaleTeam) error {
	return m.Called(tx, scaleTeam).Error(0)
}

func (m *ScaleTeamManagerMock) Get(tx *gorm.DB, options ...db.GetOption) ([]db.ScaleTeam, error) {
	toReturn := m.Called(tx, options)
	return toReturn.Get(0).([]db.ScaleTeam), toReturn.Error(1)
}

type UserManagerMock struct {
	mock.Mock
}

func (m *UserManagerMock) DB() *gorm.DB {
	return m.Called().Get(0).(*gorm.DB)
}

func (m *UserManagerMock) Create(tx *gorm.DB, scaleTeamID int, login string, status db.UserStatus) (db.User, error) {
	toReturn := m.Called(tx, scaleTeamID, login, status)
	return toReturn.Get(0).(db.User), toReturn.Error(1)
}

func (m *UserManagerMock) Update(tx *gorm.DB, scaleTeam db.User) error {
	return m.Called(tx, scaleTeam).Error(0)
}

func (m *UserManagerMock) Delete(tx *gorm.DB, scaleTeam db.User) error {
	return m.Called(tx, scaleTeam).Error(0)
}

func (m *UserManagerMock) Get(tx *gorm.DB, options ...db.GetOption) ([]db.User, error) {
	toReturn := m.Called(tx, options)
	return toReturn.Get(0).([]db.User), toReturn.Error(1)
}

type ScaleTeamMock struct {
	mock.Mock
}

func (m *ScaleTeamMock) GetID() int {
	return m.Called().Int(0)
}

func (m *ScaleTeamMock) GetBeginAt() time.Time {
	return m.Called().Get(0).(time.Time)
}

func (m *ScaleTeamMock) GetNotified() bool {
	return m.Called().Bool(0)
}

func (m *ScaleTeamMock) Get(tx *gorm.DB, options ...db.GetOption) ([]db.User, error) {
	toReturn := m.Called(tx, options)
	return toReturn.Get(0).([]db.User), toReturn.Error(1)
}

func (m *ScaleTeamMock) SetID(id int) {
	m.Called(id)
}

func (m *ScaleTeamMock) SetBeginAt(beginAt time.Time) {
	m.Called(beginAt)
}

func (m *ScaleTeamMock) SetNotified(notified bool) {
	m.Called(notified)
}

func (m *ScaleTeamMock) Save(tx *gorm.DB) error {
	return m.Called(tx).Error(0)
}

func (m *ScaleTeamMock) Delete(tx *gorm.DB) error {
	return m.Called(tx).Error(0)
}

type UserMock struct {
	mock.Mock
}

func (m *UserMock) GetID() int {
	return m.Called().Int(0)
}

func (m *UserMock) GetScaleTeamID() int {
	return m.Called().Int(0)
}

func (m *UserMock) GetLogin() string {
	return m.Called().String(0)
}

func (m *UserMock) GetStatus() db.UserStatus {
	return db.UserStatus(m.Called().String(0))
}

func (m *UserMock) SetScaleTeamID(id int) {
	m.Called(id)
}

func (m *UserMock) SetLogin(login string) {
	m.Called(login)
}

func (m *UserMock) SetStatus(status db.UserStatus) {
	m.Called(status)
}

func (m *UserMock) Save(tx *gorm.DB) error {
	return m.Called(tx).Error(0)
}

func (m *UserMock) Delete(tx *gorm.DB) error {
	return m.Called(tx).Error(0)
}
