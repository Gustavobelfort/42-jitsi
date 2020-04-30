package oauth

import (
	"net/http/httptest"

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

	m.router.Any("/*path", func(ctx *gin.Context) {
		path := ctx.Param("path")

		query := ctx.Request.URL.Query()
		body := gin.H{}
		if err := ctx.ShouldBindJSON(&body); err != nil {
			body = nil
		}

		toReturn := m.MethodCalled(path, ctx.Request.Method, query, body, ctx.Request.Header.Get("Authorization"))
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
