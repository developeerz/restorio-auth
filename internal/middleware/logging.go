package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func Logging(ctx *gin.Context) {
	bodyBytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Error().AnErr("MiddlewareLogging", err).Send()
	}

	ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	tsReq := time.Now()

	ctx.Next()

	log.Info().
		Str("method", ctx.Request.Method).
		Str("url", ctx.Request.RequestURI).
		Int("status", ctx.Writer.Status()).
		Str("body", string(bodyBytes)).
		Dur("duration", time.Since(tsReq)).
		Send()
}
