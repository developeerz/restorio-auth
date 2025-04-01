package routers

import (
	"github.com/developeerz/restorio-auth/internal/handler"
	"github.com/gin-gonic/gin"
)

func NewAuthRouter(router *gin.Engine, authHandler *handler.AuthHandler) {
	authapi := router.Group("/api/auth")

	authapi.GET("/check-access", authHandler.CheckAccess)
	authapi.GET("/refresh", authHandler.Refresh)
}
