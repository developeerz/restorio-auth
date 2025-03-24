package handler

import (
	"net/http"
	"strings"

	"github.com/developeerz/restorio-auth/internal/dto"
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
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, &dto.Error{Message: "Authorization header required"})
		return
	}

	token, msg := extractToken(authHeader)
	if msg != nil {
		ctx.JSON(http.StatusUnauthorized, msg)
		return
	}

	jwt, err := authHandler.authService.Refresh(token)
	if err != nil {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	ctx.JSON(http.StatusOK, jwt)
}

func (authHandler *AuthHandler) CheckAccess(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	if authHeader == "" {
		ctx.JSON(http.StatusUnauthorized, &dto.Error{Message: "Authorization header required"})
		return
	}

	accessToken, msg := extractToken(authHeader)
	if msg != nil {
		ctx.JSON(http.StatusUnauthorized, msg)
		return
	}

	id, roles, err := jwt.ParseAccess(accessToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, msg)
		return
	}

	ctx.Header("X-User-Id", id)
	ctx.Header("X-Roles", strings.Join(roles, ","))

	ctx.Status(http.StatusOK)
}

func extractToken(authHeader string) (string, *dto.Error) {
	if authHeader == "" {
		return "", &dto.Error{Message: "Authorization header required"}
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", &dto.Error{Message: "Invalid authorization format"}
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}
