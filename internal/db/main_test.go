package db

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestManagers(t *testing.T) {
	suite.Run(t, &ManagerSuite{})
}

type ScaleTeamManagerMock struct {
	mock.Mock
}

func (sMock *ScaleTeamManagerMock) Create(_ int, _ time.Time, _ bool) (ScaleTeam, error) {
	sMock.Called()
	return nil, nil
}

func (sMock *ScaleTeamManagerMock) CreateWithTransaction(_ *gorm.DB, _ int, _ time.Time, _ bool) (ScaleTeam, error) {
	sMock.Called()
	return nil, nil
}

func (sMock *ScaleTeamManagerMock) Get(_ ...GetOption) ([]ScaleTeam, error) {
	sMock.Called()
	return nil, nil
}

func (sMock *ScaleTeamManagerMock) GetWithTransaction(_ *gorm.DB, _ ...GetOption) ([]ScaleTeam, error) {
	sMock.Called()
	return nil, nil
}

func (sMock *ScaleTeamManagerMock) Update(scaleTeam ScaleTeam) error {
	return sMock.Called(scaleTeam).Error(0)
}

func (sMock *ScaleTeamManagerMock) UpdateWithTransaction(tx *gorm.DB, scaleTeam ScaleTeam) error {
	return sMock.Called(tx, scaleTeam).Error(0)
}

func (sMock *ScaleTeamManagerMock) Delete(scaleTeam ScaleTeam) error {
	return sMock.Called(scaleTeam).Error(0)
}

func (sMock *ScaleTeamManagerMock) DeleteWithTransaction(tx *gorm.DB, scaleTeam ScaleTeam) error {
	return sMock.Called(tx, scaleTeam).Error(0)
}

func TestScaleTeamModel(t *testing.T) {
	assert := assert.New(t)

	var (
		expectedID      = 1
		expectedBeginAt = time.Now()
		expectedNotifed = true
	)

	scaleTeam := &scaleTeamModel{}

	assert.Implements((*ScaleTeam)(nil), scaleTeam)

	scaleTeam.SetID(expectedID)
	scaleTeam.SetBeginAt(expectedBeginAt)
	scaleTeam.SetNotified(expectedNotifed)

	assert.Equal(expectedID, scaleTeam.GetID())
	assert.Equal(expectedBeginAt, scaleTeam.GetBeginAt())
	assert.Equal(expectedNotifed, scaleTeam.GetNotified())

	expectedError := errors.New("testing error")

	db, _, err := sqlmock.New()
	require.NoError(t, err)

	tx, err := gorm.Open("postgres", db)
	require.NoError(t, err)

	mock := &ScaleTeamManagerMock{}
	scaleTeam.scaleTeamManager = mock

	mock.On("Update", scaleTeam).Return(expectedError)
	mock.On("UpdateWithTransaction", tx, scaleTeam).Return(expectedError)
	mock.On("Delete", scaleTeam).Return(expectedError)
	mock.On("DeleteWithTransaction", tx, scaleTeam).Return(expectedError)

	assert.Equal(expectedError, scaleTeam.Save())
	assert.Equal(expectedError, scaleTeam.SaveWithTransaction(tx))
	assert.Equal(expectedError, scaleTeam.Delete())
	assert.Equal(expectedError, scaleTeam.DeleteWithTransaction(tx))
}

type UserManagerMock struct {
	mock.Mock
}

func (sMock *UserManagerMock) Create(_ int, _ string, _ UserStatus) (User, error) {
	sMock.Called()
	return nil, nil
}

func (sMock *UserManagerMock) CreateWithTransaction(_ *gorm.DB, _ int, _ string, _ UserStatus) (User, error) {
	sMock.Called()
	return nil, nil
}

func (sMock *UserManagerMock) Get(_ ...GetOption) ([]User, error) {
	sMock.Called()
	return nil, nil
}

func (sMock *UserManagerMock) GetWithTransaction(_ *gorm.DB, _ ...GetOption) ([]User, error) {
	sMock.Called()
	return nil, nil
}

func (sMock *UserManagerMock) Update(user User) error {
	return sMock.Called(user).Error(0)
}

func (sMock *UserManagerMock) UpdateWithTransaction(tx *gorm.DB, user User) error {
	return sMock.Called(tx, user).Error(0)
}

func (sMock *UserManagerMock) Delete(user User) error {
	return sMock.Called(user).Error(0)
}

func (sMock *UserManagerMock) DeleteWithTransaction(tx *gorm.DB, user User) error {
	return sMock.Called(tx, user).Error(0)
}

func TestUserModel(t *testing.T) {
	assert := assert.New(t)

	var (
		expectedID          = 1
		expectedScaleTeamID = 2
		expectedLogin       = "xlogin"
		expectedStatus      = Corrector
	)

	user := &userModel{
		ID: expectedID,
	}

	assert.Implements((*User)(nil), user)

	user.SetScaleTeamID(expectedScaleTeamID)
	user.SetLogin(expectedLogin)
	user.SetStatus(expectedStatus)

	assert.Equal(expectedID, user.GetID())
	assert.Equal(expectedScaleTeamID, user.GetScaleTeamID())
	assert.Equal(expectedLogin, user.GetLogin())
	assert.Equal(expectedStatus, user.GetStatus())

	expectedError := errors.New("testing error")

	mock := &UserManagerMock{}
	user.userManager = mock

	db, _, err := sqlmock.New()
	require.NoError(t, err)

	tx, err := gorm.Open("postgres", db)
	require.NoError(t, err)

	mock.On("Update", user).Return(expectedError)
	mock.On("UpdateWithTransaction", tx, user).Return(expectedError)
	mock.On("Delete", user).Return(expectedError)
	mock.On("DeleteWithTransaction", tx, user).Return(expectedError)

	assert.Equal(expectedError, user.Save())
	assert.Equal(expectedError, user.SaveWithTransaction(tx))
	assert.Equal(expectedError, user.Delete())
	assert.Equal(expectedError, user.DeleteWithTransaction(tx))
}
