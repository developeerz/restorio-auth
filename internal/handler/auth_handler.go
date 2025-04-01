package handler

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (authHandler *AuthHandler) Refresh(ctx *gin.Context) {
	refreshOld, err := ctx.Cookie("refresh")
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	access, refresh, err := authHandler.authService.Refresh(refreshOld)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	ctx.SetCookie("refresh", refresh, jwt.RefreshMaxAge, "/api/auth/refresh", "", false, true)
	ctx.JSON(http.StatusOK, access)
}

func (authHandler *AuthHandler) CheckAccess(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	accessToken, err := extractToken(authHeader)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id, roles, err := jwt.GetAccess(accessToken)
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx.Header("X-User-Id", id)
	ctx.Header("X-Roles", strings.Join(roles, ","))

	ctx.Status(http.StatusOK)
}

func extractToken(authHeader string) (string, error) {
	if authHeader == "" {
		return "", fmt.Errorf("Empty auth header")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("Invalid authorization format")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}
