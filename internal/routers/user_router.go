package routers

import (
	user_handler "github.com/developeerz/restorio-auth/internal/handler/user"
	"github.com/gin-gonic/gin"
)

const (
	userGroupPath       = "/api/auth-service/user"
	userGroupSignUpPath = "/sign-up"
	userGroupLoginPath  = "/login"
	userGroupVerifyPath = "/verify"
)

func NewUserRouter(router *gin.Engine, userHandler *user_handler.Handler) {
	userapi := router.Group(userGroupPath)

	userapi.POST(userGroupSignUpPath, userHandler.SignUp)
	userapi.POST(userGroupLoginPath, userHandler.Login)
	userapi.POST(userGroupVerifyPath, userHandler.Verification)
}
