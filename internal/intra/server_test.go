package intra

import (
	"net/http/httptest"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type ServerMock struct {
	mock.Mock

	router *gin.Engine
	Server *httptest.Server
}

func (m *ServerMock) initRouter() {
	m.router = gin.New()

	// Mocking oauth requests
	// Assuming that oauth lib works
	m.router.POST("/oauth/token", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"access_token": "token",
			"token_type":   "Bearer",
			"created_at":   time.Now().Unix(),
			"expires_in":   7200,
		})
	})

	// Mocking user show request
	m.router.GET("/v2/users/:id", func(ctx *gin.Context) {
		toReturn := m.MethodCalled("GetUser", ctx.Param("id"))
		for key, value := range toReturn.Get(2).(gin.H) {
			ctx.Header(key, value.(string))
		}
		ctx.JSON(toReturn.Int(0), toReturn.Get(1))
	})

	// Mocking team show users index request
	m.router.GET("/v2/teams/:id/users", func(ctx *gin.Context) {
		toReturn := m.MethodCalled("GetTeamUsers", ctx.Param("id"))
		for key, value := range toReturn.Get(2).(gin.H) {
			ctx.Header(key, value.(string))
		}
		ctx.JSON(toReturn.Int(0), toReturn.Get(1))
	})
}

func NewServerMock() *ServerMock {
	mock := &ServerMock{}
	mock.initRouter()
	mock.Server = httptest.NewServer(mock.router)
	return mock
}
