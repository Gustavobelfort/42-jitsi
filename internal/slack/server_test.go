package slack

import (
	"encoding/json"
	"net/http"
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

	m.router.POST("/", func(ctx *gin.Context) {
		if !validateRequest(ctx.Request) {
			ctx.JSON(500, gin.H{})
		}
		ctx.JSON(201, gin.H{})
	})
}

func validateRequest(r *http.Request) bool {

	var p PostMessageParameters
	json.NewDecoder(r.Body).Decode(&p)
	if p.Workspace != "testWorkspace" || p.Attachments[0].TitleLink != "https://meet.jit.si/1-xlogin" {
		return false
	}
	return true
}

func NewServerMock() *ServerMock {
	mock := &ServerMock{}
	mock.initRouter()
	mock.Server = httptest.NewServer(mock.router)
	return mock
}
