package routers

import (
	auth_handler "github.com/developeerz/restorio-auth/internal/handler/auth"
	"github.com/gin-gonic/gin"
)

const (
	authGroupPath            = "/api/auth-service/auth"
	authGroupCheckAccessPath = "/check-access"
	authGroupRefreshPath     = "/refresh"

	AuthGroupFullRefreshPath = authGroupPath + authGroupRefreshPath
)

func NewAuthRouter(router *gin.Engine, authHandler *auth_handler.Handler) {
	authapi := router.Group(authGroupPath)

	authapi.GET(authGroupCheckAccessPath, authHandler.CheckAccess)
	authapi.GET(authGroupRefreshPath, authHandler.Refresh)
}
