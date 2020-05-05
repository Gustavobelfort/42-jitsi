package router

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/sirupsen/logrus"
)

var (
	// AlreadyStartedError is returned when you try to start a router that was already started.
	AlreadyStartedError = errors.New("the router was already started")
)

func handleError(ctx *gin.Context, err error) {
	logError := &logging.WithLogError{}
	if errors.As(err, &logError) {
		if logError.LogLevel <= logrus.WarnLevel {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Bad request.", "details": err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error.", "details": nil})
}
