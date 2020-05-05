package router

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/sirupsen/logrus"
)

func (r *Router) setupEngine(prefix string) {
	r.engine = gin.New()
	group := r.engine.Group(prefix)

	r.setupMiddlewares(group)

	// Webhook receiving part. Thanks to the middleware it can assume that the received headers are valid.
	// However some precautions were taken just in case.
	//
	// The context should be populated for log too.
	group.POST("/webhooks", func(ctx *gin.Context) {
		rCtx := ctx.Request.Context()
		ctxlogger := logging.ContextLog(rCtx, logrus.StandardLogger())

		body, err := ioutil.ReadAll(ctx.Request.Body)
		if err != nil {
			ctxlogger.Errorf("while reading the request's body: %v", err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error.", "details": nil})
			return
		}

		switch ctx.Request.Header.Get("X-Event") {
		case "create":
			err = r.handler.HandleCreate(rCtx, body)
		case "update":
			err = r.handler.HandleUpdate(rCtx, body)
		case "destroy":
			err = r.handler.HandleDestroy(rCtx, body)
		default:
			err = errors.New("unknown event")
		}
		logging.LogError(ctxlogger, err, "while handling request")
		if err != nil {
			handleError(ctx, err)
			return
		}

		ctx.Status(http.StatusNoContent)
	})
}
