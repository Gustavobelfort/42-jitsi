package tasks

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/go-playground/assert.v1"
)

func TestScaleTeamHandler(t *testing.T) {
	t.Run("NewTasksHandler", func(t *testing.T) {
		client := &ClientMock{}
		db := &gorm.DB{}

		handler := NewTasksHandler(client, db)
		require.IsType(t, &tasksHandler{}, handler)

		tHandler := handler.(*tasksHandler)

		assert.Equal(t, db, tHandler.db)
		assert.Equal(t, db, tHandler.scaleTeamManager.DB())
		assert.Equal(t, db, tHandler.userManager.DB())
		assert.Equal(t, client, tHandler.client)
	})

	suite.Run(t, new(TasksHandlerSuite))
}

type TasksHandlerSuite struct {
	suite.Suite

	handler *tasksHandler

	stMock *ScaleTeamManagerMock
	uMock  *UserManagerMock
	cMock  *ClientMock

	db     *gorm.DB
	dbMock sqlmock.Sqlmock
}
