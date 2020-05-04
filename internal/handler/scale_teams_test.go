package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gustavobelfort/42-jitsi/internal/db"
	"github.com/jinzhu/gorm"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestScaleTeamHandler(t *testing.T) {
	t.Run("NewScaleTeamHander", func(t *testing.T) {
		client := &ClientMock{}
		db := &gorm.DB{}

		handler := NewScaleTeamHandler(client, db)
		require.IsType(t, &scaleTeamHandler{}, handler)

		stHandler := handler.(*scaleTeamHandler)

		assert.Equal(t, db, stHandler.db)
		assert.Equal(t, db, stHandler.scaleTeamManager.DB())
		assert.Equal(t, db, stHandler.userManager.DB())
		assert.Equal(t, client, stHandler.client)
	})

	suite.Run(t, new(ScaleTeamHandlerSuite))
}

type ScaleTeamHandlerSuite struct {
	suite.Suite

	handler *scaleTeamHandler

	stMock *ScaleTeamManagerMock
	uMock  *UserManagerMock
	cMock  *ClientMock

	db     *gorm.DB
	dbMock sqlmock.Sqlmock
}

func (s *ScaleTeamHandlerSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.dbMock, err = sqlmock.New()
	s.Require().NoError(err)

	s.db, err = gorm.Open("postgres", db)
	s.Require().NoError(err)
}

func (s *ScaleTeamHandlerSuite) SetupTest() {
	s.stMock = &ScaleTeamManagerMock{}
	s.uMock = &UserManagerMock{}
	s.cMock = &ClientMock{}

	s.handler = &scaleTeamHandler{
		db:               s.db,
		scaleTeamManager: s.stMock,
		userManager:      s.uMock,

		client: s.cMock,
	}
}

func (s *ScaleTeamHandlerSuite) Test00_HandleCreate() {
	expectedID := 21
	expectedCorrector := "xlogin"
	expectedTeam := 42

	payload := []byte(fmt.Sprintf(
		`{"id": %d, "user": {"login": "%s"}, "team": {"id": %d}, "begin_at": "2020-07-15T21:00:00.000Z"}`,
		expectedID,
		expectedCorrector,
		expectedTeam,
	))

	expectedLogins := []string{"ylogin"}

	expectedContext := context.Background()
	s.cMock.On("GetTeamMembers", expectedContext, expectedTeam).Return(expectedLogins, nil).Once()

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()

	recordMock := &ScaleTeamMock{}
	defer recordMock.AssertExpectations(s.T())
	s.stMock.On("Create", mock.Anything, expectedID, mock.Anything, false).Return(recordMock, nil).Once()

	s.uMock.On("Create", mock.Anything, expectedID, expectedCorrector, db.Corrector).Return(&UserMock{}, nil).Once()
	s.uMock.On("Create", mock.Anything, expectedID, expectedLogins[0], db.Corrected).Return(&UserMock{}, nil).Once()

	recordMock.On("GetID").Return(expectedID).Twice()

	err := s.handler.HandleCreate(expectedContext, payload)
	s.NoError(err)
}

func (s *ScaleTeamHandlerSuite) Test01_HandleCreate_PayloadError() {
	payload := []byte(`{"id": "twenty-one"}`)

	err := s.handler.HandleCreate(context.Background(), payload)
	s.Error(err)
}

func (s *ScaleTeamHandlerSuite) Test02_HandleCreate_ClientError() {
	expectedTeam := 42

	payload := []byte(fmt.Sprintf(
		`{"id": 21, "user": {"login": "xlogin"}, "team": {"id": %d}, "begin_at": "2020-07-15T21:00:00.000Z"}`,
		expectedTeam,
	))

	expectedError := errors.New("testing")

	expectedContext := context.Background()
	s.cMock.On("GetTeamMembers", expectedContext, expectedTeam).Return([]string{}, expectedError).Once()

	err := s.handler.HandleCreate(expectedContext, payload)
	s.Error(err)
	s.Equal(expectedError, err)
}

func (s *ScaleTeamHandlerSuite) Test03_HandleCreate_CreateScaleTeamError() {
	expectedID := 21
	expectedCorrector := "xlogin"
	expectedTeam := 42

	payload := []byte(fmt.Sprintf(
		`{"id": %d, "user": {"login": "%s"}, "team": {"id": %d}, "begin_at": "2020-07-15T21:00:00.000Z"}`,
		expectedID,
		expectedCorrector,
		expectedTeam,
	))

	expectedLogins := []string{"ylogin"}

	expectedContext := context.Background()
	s.cMock.On("GetTeamMembers", expectedContext, expectedTeam).Return(expectedLogins, nil).Once()

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()

	expectedError := errors.New("testing")

	s.stMock.On("Create", mock.Anything, expectedID, mock.Anything, false).Return(&ScaleTeamMock{}, expectedError).Once()

	err := s.handler.HandleCreate(expectedContext, payload)
	s.Error(err)
	s.Equal(expectedError, err)
}

func (s *ScaleTeamHandlerSuite) Test04_HandleCreate_CreateCorrectorError() {
	expectedID := 21
	expectedCorrector := "xlogin"
	expectedTeam := 42

	payload := []byte(fmt.Sprintf(
		`{"id": %d, "user": {"login": "%s"}, "team": {"id": %d}, "begin_at": "2020-07-15T21:00:00.000Z"}`,
		expectedID,
		expectedCorrector,
		expectedTeam,
	))

	expectedLogins := []string{"ylogin"}

	expectedContext := context.Background()
	s.cMock.On("GetTeamMembers", expectedContext, expectedTeam).Return(expectedLogins, nil).Once()

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()

	recordMock := &ScaleTeamMock{}
	defer recordMock.AssertExpectations(s.T())
	s.stMock.On("Create", mock.Anything, expectedID, mock.Anything, false).Return(recordMock, nil).Once()

	expectedError := errors.New("testing")

	s.uMock.On("Create", mock.Anything, expectedID, expectedCorrector, db.Corrector).Return(&UserMock{}, expectedError).Once()

	recordMock.On("GetID").Return(expectedID).Once()

	err := s.handler.HandleCreate(expectedContext, payload)
	s.Error(err)
	s.Equal(expectedError, err)
}

func (s *ScaleTeamHandlerSuite) Test05_HandleCreate_CreateCorrectedError() {
	expectedID := 21
	expectedCorrector := "xlogin"
	expectedTeam := 42

	payload := []byte(fmt.Sprintf(
		`{"id": %d, "user": {"login": "%s"}, "team": {"id": %d}, "begin_at": "2020-07-15T21:00:00.000Z"}`,
		expectedID,
		expectedCorrector,
		expectedTeam,
	))

	expectedLogins := []string{"ylogin"}

	expectedContext := context.Background()
	s.cMock.On("GetTeamMembers", expectedContext, expectedTeam).Return(expectedLogins, nil).Once()

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()

	recordMock := &ScaleTeamMock{}
	defer recordMock.AssertExpectations(s.T())
	s.stMock.On("Create", mock.Anything, expectedID, mock.Anything, false).Return(recordMock, nil).Once()

	expectedError := errors.New("testing")

	s.uMock.On("Create", mock.Anything, expectedID, expectedCorrector, db.Corrector).Return(&UserMock{}, nil).Once()
	s.uMock.On("Create", mock.Anything, expectedID, expectedLogins[0], db.Corrected).Return(&UserMock{}, expectedError).Once()

	recordMock.On("GetID").Return(expectedID).Twice()

	err := s.handler.HandleCreate(expectedContext, payload)
	s.Error(err)
	s.Equal(expectedError, err)
}

func (s *ScaleTeamHandlerSuite) Test06_HandleUpdate() {
	expectedTeam := 42

	payload := []byte(fmt.Sprintf(
		`{"id": 21, "user": {"login": "xlogin"}, "team": {"id": %d}, "begin_at": "2020-07-15T21:00:00.000Z"}`,
		expectedTeam,
	))

	expectedContext := context.Background()
	s.cMock.On("GetTeamMembers", expectedContext, expectedTeam).Return([]string{}, nil).Once()

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()

	recordMock := &ScaleTeamMock{}
	defer recordMock.AssertExpectations(s.T())
	s.stMock.On("Get", mock.Anything, mock.Anything).Return([]db.ScaleTeam{recordMock}, nil).Once()

	recordMock.On("SetBeginAt", mock.Anything).Return().Once()
	recordMock.On("SetNotified", false).Return().Once()

	recordMock.On("Save", mock.Anything).Return(nil).Once()

	err := s.handler.HandleUpdate(expectedContext, payload)
	s.NoError(err)
}

func (s *ScaleTeamHandlerSuite) Test07_HandleUpdate_NotInDB() {
	expectedID := 21
	expectedCorrector := "xlogin"
	expectedTeam := 42

	payload := []byte(fmt.Sprintf(
		`{"id": %d, "user": {"login": "%s"}, "team": {"id": %d}, "begin_at": "2020-07-15T21:00:00.000Z"}`,
		expectedID,
		expectedCorrector,
		expectedTeam,
	))

	expectedLogins := []string{"ylogin"}

	expectedContext := context.Background()
	s.cMock.On("GetTeamMembers", expectedContext, expectedTeam).Return(expectedLogins, nil).Once()

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()

	s.stMock.On("Get", mock.Anything, mock.Anything).Return([]db.ScaleTeam{}, nil)

	recordMock := &ScaleTeamMock{}
	defer recordMock.AssertExpectations(s.T())
	s.stMock.On("Create", mock.Anything, expectedID, mock.Anything, false).Return(recordMock, nil).Once()

	s.uMock.On("Create", mock.Anything, expectedID, expectedCorrector, db.Corrector).Return(&UserMock{}, nil).Once()
	s.uMock.On("Create", mock.Anything, expectedID, expectedLogins[0], db.Corrected).Return(&UserMock{}, nil).Once()

	recordMock.On("GetID").Return(expectedID).Twice()

	err := s.handler.HandleUpdate(expectedContext, payload)
	s.NoError(err)
}

func (s *ScaleTeamHandlerSuite) Test08_HandleUpdate_PayloadError() {
	payload := []byte(`{"id": "twenty-one"}`)

	err := s.handler.HandleUpdate(context.Background(), payload)
	s.Error(err)
}

func (s *ScaleTeamHandlerSuite) Test09_HandleUpdate_GetError() {
	expectedTeam := 42

	payload := []byte(fmt.Sprintf(
		`{"id": 21, "user": {"login": "xlogin"}, "team": {"id": %d}, "begin_at": "2020-07-15T21:00:00.000Z"}`,
		expectedTeam,
	))

	expectedContext := context.Background()
	s.cMock.On("GetTeamMembers", expectedContext, expectedTeam).Return([]string{}, nil).Once()

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()

	expectedError := errors.New("if you read this, listen to Salut C'est Cool")
	s.stMock.On("Get", mock.Anything, mock.Anything).Return([]db.ScaleTeam{}, expectedError).Once()

	err := s.handler.HandleUpdate(expectedContext, payload)
	s.Error(err)
	s.Equal(expectedError, err)
}

func (s *ScaleTeamHandlerSuite) Test10_HandleUpdate_SaveError() {
	expectedTeam := 42

	payload := []byte(fmt.Sprintf(
		`{"id": 21, "user": {"login": "xlogin"}, "team": {"id": %d}, "begin_at": "2020-07-15T21:00:00.000Z"}`,
		expectedTeam,
	))

	expectedContext := context.Background()
	s.cMock.On("GetTeamMembers", expectedContext, expectedTeam).Return([]string{}, nil).Once()

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()

	recordMock := &ScaleTeamMock{}
	defer recordMock.AssertExpectations(s.T())
	s.stMock.On("Get", mock.Anything, mock.Anything).Return([]db.ScaleTeam{recordMock}, nil).Once()

	recordMock.On("SetBeginAt", mock.Anything).Return().Once()
	recordMock.On("SetNotified", false).Return().Once()

	expectedError := errors.New("testing")
	recordMock.On("Save", mock.Anything).Return(expectedError).Once()

	err := s.handler.HandleUpdate(expectedContext, payload)
	s.Error(err)
	s.Equal(expectedError, err)
}

func (s *ScaleTeamHandlerSuite) Test11_HandleDestroy() {
	expectedID := 21

	payload := []byte(fmt.Sprintf(
		`{"id": %d}`,
		expectedID,
	))

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectCommit()

	recordMock := &ScaleTeamMock{}
	defer recordMock.AssertExpectations(s.T())
	s.stMock.On("Get", mock.Anything, mock.Anything).Return([]db.ScaleTeam{recordMock}, nil).Once()

	recordMock.On("Delete", mock.Anything).Return(nil).Once()

	err := s.handler.HandleDestroy(context.Background(), payload)
	s.NoError(err)
}

func (s *ScaleTeamHandlerSuite) Test12_HandleDestroy_NotInDB() {
	expectedID := 21

	payload := []byte(fmt.Sprintf(
		`{"id": %d}`,
		expectedID,
	))

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()

	s.stMock.On("Get", mock.Anything, mock.Anything).Return([]db.ScaleTeam{}, nil).Once()

	err := s.handler.HandleDestroy(context.Background(), payload)
	s.Error(err)
	s.True(errors.Is(err, NotInDBError))
}

func (s *ScaleTeamHandlerSuite) Test13_HandleDestroy_PayloadError() {
	expectedID := 21

	payload := []byte(fmt.Sprintf(
		`{"idd": %d}`,
		expectedID,
	))

	err := s.handler.HandleDestroy(context.Background(), payload)
	s.Error(err)
}

func (s *ScaleTeamHandlerSuite) Test14_HandleDestroy_InvalidJSON() {
	expectedID := 21

	payload := []byte(fmt.Sprintf(
		`{"idd": %d`,
		expectedID,
	))

	err := s.handler.HandleDestroy(context.Background(), payload)
	s.Error(err)
}

func (s *ScaleTeamHandlerSuite) Test15_HandleDestroy_GetError() {
	expectedID := 21

	payload := []byte(fmt.Sprintf(
		`{"id": %d}`,
		expectedID,
	))

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()

	expectedError := errors.New("testing")
	s.stMock.On("Get", mock.Anything, mock.Anything).Return([]db.ScaleTeam{}, expectedError).Once()

	err := s.handler.HandleDestroy(context.Background(), payload)
	s.Error(err)
	s.Equal(expectedError, err)
}

func (s *ScaleTeamHandlerSuite) Test16_HandleDestroy_DeleteError() {
	expectedID := 21

	payload := []byte(fmt.Sprintf(
		`{"id": %d}`,
		expectedID,
	))

	s.dbMock.ExpectBegin()
	s.dbMock.ExpectRollback()

	recordMock := &ScaleTeamMock{}
	defer recordMock.AssertExpectations(s.T())
	s.stMock.On("Get", mock.Anything, mock.Anything).Return([]db.ScaleTeam{recordMock}, nil).Once()

	expectedError := errors.New("testing")
	recordMock.On("Delete", mock.Anything).Return(expectedError).Once()

	err := s.handler.HandleDestroy(context.Background(), payload)
	s.Error(err)
	s.Equal(expectedError, err)
}

func (s *ScaleTeamHandlerSuite) TearDownTest() {
	s.stMock.AssertExpectations(s.T())
	s.uMock.AssertExpectations(s.T())
	s.cMock.AssertExpectations(s.T())
	s.NoError(s.dbMock.ExpectationsWereMet())
}
