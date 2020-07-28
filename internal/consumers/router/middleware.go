package router

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gustavobelfort/42-jitsi/internal/logging"
	"github.com/sirupsen/logrus"
)

func checkToken(model, event, secret string, registries map[string]string) bool {
	expected, ok := registries[fmt.Sprintf("%s.%s", model, event)]
	return ok && expected == secret
}

// contextMiddleware populates the request's context with the global context and logging values.
func (r *Router) contextMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		model, event, delivery := ctx.GetHeader("X-Model"), ctx.GetHeader("X-Event"), ctx.GetHeader("X-Delivery")

		var remoteAddr string
		if remoteAddr = ctx.GetHeader("X-Forwarded-For"); remoteAddr == "" {
			remoteAddr = ctx.Request.RemoteAddr
		}

		rCtx := logging.ContextWithFields(r.ctx, logrus.Fields{
			"model":       model,
			"event":       event,
			"delivery_id": delivery,
			"remote_addr": remoteAddr,
		})
		rCtx, cancel := context.WithTimeout(ctx, r.timeout)
		defer cancel()
		ctx.Request = ctx.Request.WithContext(rCtx)
		ctx.Next()
	}
}

// recoverMiddleware defers a recover function that will log any panic error that occurs and send an internal error
// signal.
func (r *Router) recoverMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ctxlogger := logging.ContextLog(ctx.Request.Context(), logrus.StandardLogger())
				ctxlogger.WithField("error", err).Errorf("request paniqued: %v", err)
				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error.", "details": nil})
			}
		}()
		ctx.Next()
	}
}

// validationMiddleware verifies that the request is authenticated and that the model is supported.
func (r *Router) validationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxlogger := logging.ContextLog(ctx.Request.Context(), logrus.StandardLogger())

		model, event, secret := ctx.GetHeader("X-Model"), ctx.GetHeader("X-Event"), ctx.GetHeader("X-Secret")

		if !checkToken(model, event, secret, r.registries) {
			ctxlogger.Warn("unauthorized request: denying request")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized."})
			return
		}

		if model != "scale_team" {
			ctxlogger.Warn("unhandled model: denying request")
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Unhandled model."})
			return
		}

		switch event {
		case "create", "update", "destroy":
			break
		default:
			ctxlogger.Warn("unhandled event: denying request")
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Unhandled event."})
			return
		}
	}
}

func (r *Router) setupMiddlewares(group *gin.RouterGroup) {
	group.Use(
		r.recoverMiddleware(),
		r.contextMiddleware(),
		r.validationMiddleware(),
	)
}
