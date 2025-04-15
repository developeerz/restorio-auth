package auth

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

const (
	cookieRefreshName = "refresh"
	headerUserIDKey   = "X-Telegram-Id"
	headerRolesKey    = "X-Roles"
)

type Handler struct {
	service     Service
	refreshPath string
}

func NewHandler(service Service, refreshPath string) *Handler {
	return &Handler{service: service, refreshPath: refreshPath}
}

func (handler *Handler) Refresh(ctx *gin.Context) {
	refreshOld, err := ctx.Cookie(cookieRefreshName)
	if err != nil {
		log.Error().AnErr("Refresh", err).Send()
		ctx.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	access, refresh, err := handler.service.Refresh(refreshOld)
	if err != nil {
		log.Error().AnErr("Refresh", err).Send()
		ctx.Status(http.StatusUnauthorized)

		return
	}

	ctx.SetCookie(cookieRefreshName, refresh, jwt.RefreshMaxAge, handler.refreshPath, "", false, true)
	ctx.JSON(http.StatusOK, access)
}

func (handler *Handler) CheckAccess(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		log.Error().AnErr("Refresh", fmt.Errorf("cannot find \"Authorization\" header"))
		ctx.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	accessToken, err := extractToken(authHeader)
	if err != nil {
		log.Error().AnErr("Refresh", err).Send()
		ctx.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	id, roles, err := jwt.GetAccess(accessToken)
	if err != nil {
		log.Error().AnErr("Refresh", err).Send()
		ctx.AbortWithStatus(http.StatusUnauthorized)

		return
	}

	ctx.Header(headerUserIDKey, id)
	ctx.Header(headerRolesKey, strings.Join(roles, ","))

	ctx.Status(http.StatusOK)
}

func extractToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("empty auth header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("invalid authorization format")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}
