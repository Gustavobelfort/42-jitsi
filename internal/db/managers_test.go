package db

import (
	"database/sql"
	"regexp"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/suite"
)

type ManagerSuite struct {
	suite.Suite
	db   *gorm.DB
	mock sqlmock.Sqlmock

	scaleTeamManager *scaleTeamManager
	scaleTeam        *scaleTeamModel

	userManager *userManager
	corrected   *userModel
	corrector   *userModel
}

/*
 * It seems like gorm and sqlmock are incompatible on UPDATE and DELETE requests. Although everything matches, the comparison fails.
 * Will implement those tests later.
 */

func (s *ManagerSuite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	// You can't run the tests properly if those interfaces are not implemented.
	s.Require().Implements((*ScaleTeamManager)(nil), &scaleTeamManager{})
	s.Require().Implements((*ScaleTeam)(nil), &scaleTeamModel{})
	s.Require().Implements((*UserManager)(nil), &userManager{})
	s.Require().Implements((*User)(nil), &userModel{})

	db, s.mock, err = sqlmock.New()
	s.Require().NoError(err)

	s.db, err = gorm.Open("postgres", db)
	s.Require().NoError(err)

	s.scaleTeamManager = &scaleTeamManager{db: s.db}
	s.userManager = &userManager{db: s.db}

	s.db.LogMode(true)
}

func (s *ManagerSuite) AfterTest() {
	s.Require().NoError(s.mock.ExpectationsWereMet())
}

func (s *ManagerSuite) Test00_CreateScaleTeam() {
	var (
		expectedID       = 1
		expectedBeginAt  = time.Now()
		expectedNotified = true
	)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "scale_teams" ("id","begin_at","notified") VALUES ($1,$2,$3) RETURNING "scale_teams"."id"`),
	).
		WithArgs(expectedID, expectedBeginAt, expectedNotified).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(expectedID))
	s.mock.ExpectCommit()

	scaleTeam, err := s.scaleTeamManager.Create(s.db, expectedID, expectedBeginAt, expectedNotified)
	s.Require().NoError(err)
	s.Require().NotNil(scaleTeam)

	s.scaleTeam = scaleTeam.(*scaleTeamModel)
}

func (s *ManagerSuite) Test01_SelectScaleTeams() {
	if s.scaleTeam == nil {
		s.T().SkipNow()
	}

	var (
		expectedID       = s.scaleTeam.ID
		expectedBeginAt  = s.scaleTeam.BeginAt
		expectedNotified = s.scaleTeam.Notified
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "scale_teams"`)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "begin_at", "notified"}).AddRow(expectedID, expectedBeginAt, expectedNotified),
		)

	scaleTeams, err := s.scaleTeamManager.Get(s.db)
	s.Require().NoError(err)
	s.Require().NotNil(scaleTeams)
	s.Require().Len(scaleTeams, 1)

	s.Assert().Equal(s.scaleTeam, scaleTeams[0])
}

func (s *ManagerSuite) Test02_SelectScaleTeamsWithOptions_0() {
	var (
		expectedID       = 1000
		expectedBeginAt  = time.Now()
		expectedNotified = false
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "scale_teams" WHERE (id = $1) AND (begin_at >= $2) AND ("scale_teams"."notified" = $3)`)).
		WithArgs(expectedID, expectedBeginAt.Format(time.RFC3339), expectedNotified).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "begin_at", "notified"}),
		)

	scaleTeams, err := s.scaleTeamManager.Get(s.db, ScaleTeamIDOption(expectedID), ScaleTeamBeginAtAfterOption(expectedBeginAt), ScaleTeamNotifiedOption(expectedNotified))
	s.Require().NoError(err)
	s.Require().NotNil(scaleTeams)
	s.Require().Len(scaleTeams, 0)
}

func (s *ManagerSuite) Test03_SelectScaleTeamsWithOptions_1() {
	var (
		expectedID       = 1000
		expectedBeginAt  = time.Now()
		expectedNotified = false
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(`WHERE (id = $1) AND (begin_at <= $2) AND ("scale_teams"."notified" = $3)`)).
		WithArgs(expectedID, expectedBeginAt.Format(time.RFC3339), expectedNotified).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "begin_at", "notified"}),
		)

	scaleTeams, err := s.scaleTeamManager.Get(s.db, ScaleTeamIDOption(expectedID), ScaleTeamBeginAtBeforeOption(expectedBeginAt), ScaleTeamNotifiedOption(expectedNotified))
	s.Require().NoError(err)
	s.Require().NotNil(scaleTeams)
	s.Require().Len(scaleTeams, 0)
}

func (s *ManagerSuite) Test04_SelectScaleTeamsWithOptions_2() {
	var (
		expectedID = 1000
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "scale_teams" WHERE (id = $1)`)).
		WithArgs(expectedID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "begin_at", "notified"}),
		)

	scaleTeams, err := s.scaleTeamManager.Get(s.db, ScaleTeamIDOption(expectedID))
	s.Require().NoError(err)
	s.Require().NotNil(scaleTeams)
	s.Require().Len(scaleTeams, 0)
}

func (s *ManagerSuite) Test05_UpdateScaleTeam() {
	s.T().Skip("UPDATE and DELETE requests are not recognized by sqlmock.")
	s.T().SkipNow()
}

func (s *ManagerSuite) Test06_CreateUserCorrector() {
	var (
		expectedID     = 2
		expectedLogin  = "xlogin"
		expectedStatus = Corrector
	)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "users" ("scale_team_id","login","status") VALUES ($1,$2,$3) RETURNING "users"."id"`),
	).
		WithArgs(expectedID, expectedLogin, expectedStatus).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	user, err := s.userManager.Create(s.db, expectedID, expectedLogin, expectedStatus)
	s.Require().NoError(err)
	s.Require().NotNil(user)

	s.corrector = user.(*userModel)
}

func (s *ManagerSuite) Test07_CreateUserCorrected() {
	var (
		expectedID     = 2
		expectedLogin  = "ylogin"
		expectedStatus = Corrected
	)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(
		regexp.QuoteMeta(`INSERT INTO "users" ("scale_team_id","login","status") VALUES ($1,$2,$3) RETURNING "users"."id"`),
	).
		WithArgs(expectedID, expectedLogin, expectedStatus).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(2))
	s.mock.ExpectCommit()

	user, err := s.userManager.Create(s.db, expectedID, expectedLogin, expectedStatus)
	s.Require().NoError(err)
	s.Require().NotNil(user)

	s.corrected = user.(*userModel)
}

func (s *ManagerSuite) Test08_SelectUsers() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users"`)).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "scale_team_id", "login", "status"}).
				AddRow(s.corrector.ID, s.corrector.ScaleTeamID, s.corrector.Login, s.corrector.Status).
				AddRow(s.corrected.ID, s.corrected.ScaleTeamID, s.corrected.Login, s.corrected.Status),
		)

	users, err := s.userManager.Get(s.db)
	s.Require().NoError(err)
	s.Require().NotNil(users)
	s.Require().Len(users, 2)
}

func (s *ManagerSuite) Test09_SelectUsersWithOptions_0() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`status = $1`)).
		WithArgs(s.corrector.Status).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "scale_team_id", "login", "status"}).
				AddRow(s.corrector.ID, s.corrector.ScaleTeamID, s.corrector.Login, s.corrector.Status),
		)

	users, err := s.userManager.Get(s.db, UserStatusOption(s.corrector.Status))
	s.Require().NoError(err)
	s.Require().NotNil(users)
	s.Require().Len(users, 1)
}

func (s *ManagerSuite) Test10_SelectUsersWithOptions_1() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`login = $1`)).
		WithArgs(s.corrected.Login).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "scale_team_id", "login", "status"}).
				AddRow(s.corrected.ID, s.corrected.ScaleTeamID, s.corrected.Login, s.corrected.Status),
		)

	users, err := s.userManager.Get(s.db, UserLoginOption(s.corrected.Login))
	s.Require().NoError(err)
	s.Require().NotNil(users)
	s.Require().Len(users, 1)
}

func (s *ManagerSuite) Test11_SelectUsersWithOptions_2() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`scale_team_id = $1`)).
		WithArgs(s.scaleTeam.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "scale_team_id", "login", "status"}).
				AddRow(s.corrector.ID, s.corrector.ScaleTeamID, s.corrector.Login, s.corrector.Status).
				AddRow(s.corrected.ID, s.corrected.ScaleTeamID, s.corrected.Login, s.corrected.Status),
		)

	users, err := s.userManager.Get(s.db, UserScaleTeamOption(s.scaleTeam.ID))
	s.Require().NoError(err)
	s.Require().NotNil(users)
	s.Require().Len(users, 2)
}

func (s *ManagerSuite) Test12_SelectUsersWithOptions_3() {
	s.mock.ExpectQuery(regexp.QuoteMeta(`id = $1`)).
		WithArgs(s.corrector.ID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "scale_team_id", "login", "status"}).
				AddRow(s.corrector.ID, s.corrector.ScaleTeamID, s.corrector.Login, s.corrector.Status),
		)

	users, err := s.userManager.Get(s.db, UserIDOption(s.corrector.ID))
	s.Require().NoError(err)
	s.Require().NotNil(users)
	s.Require().Len(users, 1)
}

func (s *ManagerSuite) Test13_UpdateUser() {
	s.T().Skip("UPDATE and DELETE requests are not recognized by sqlmock.")
	s.T().SkipNow()
}

func (s *ManagerSuite) Test14_DeleteUser() {
	s.T().Skip("UPDATE and DELETE requests are not recognized by sqlmock.")
	s.T().SkipNow()
}

func (s *ManagerSuite) Test15_DeleteScaleTeam() {
	s.T().Skip("UPDATE and DELETE requests are not recognized by sqlmock.")
	s.T().SkipNow()
}

func (s *ManagerSuite) Test16_ScaleTeamErrorCases() {
	scaleTeam, err := s.scaleTeamManager.Create(s.db, 1, time.Now(), false)
	s.Error(err)
	s.Nil(scaleTeam)

	scaleTeams, err := s.scaleTeamManager.Get(s.db)
	s.Error(err)
	s.Nil(scaleTeams)

	s.Error(s.scaleTeamManager.Update(s.db, s.scaleTeam))
	s.Error(s.scaleTeamManager.Delete(s.db, s.scaleTeam))
}

func (s *ManagerSuite) Test17_UserErrorCases() {
	user, err := s.userManager.Create(s.db, 1, "zlogin", Corrected)
	s.Error(err)
	s.Nil(user)

	users, err := s.userManager.Get(s.db)
	s.Error(err)
	s.Nil(users)

	s.Error(s.userManager.Update(s.db, s.corrected))
	s.Error(s.userManager.Delete(s.db, s.corrector))
}
