package routers

import (
	"github.com/developeerz/restorio-auth/internal/handler"
	"github.com/gin-gonic/gin"
)

func NewUserRouter(router *gin.Engine, userHandler *handler.UserHandler) {
	userapi := router.Group("/api/user")

	userapi.POST("/sign-up", userHandler.SignUp)
	userapi.POST("/login", userHandler.Login)
	userapi.POST("/verify", userHandler.Verification)
}
