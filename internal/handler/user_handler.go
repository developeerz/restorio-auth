package handler

import (
	"net/http"

	"github.com/developeerz/restorio-auth/internal/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/developeerz/restorio-auth/internal/service/user"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService user.UserService
}

func NewUserHandler(userService *user.UserService) *UserHandler {
	return &UserHandler{userService: *userService}
}

func (handler *UserHandler) SignUp(ctx *gin.Context) {
	var err error
	var req dto.SignUpRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	status, err := handler.userService.SignUp(&req)
	if err != nil || status != http.StatusOK {
		ctx.AbortWithStatus(status)
		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *UserHandler) Verification(ctx *gin.Context) {
	var err error
	var req dto.VerificationRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	status, err := handler.userService.Verify(&req)
	if err != nil {
		ctx.AbortWithStatus(status)
		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *UserHandler) Login(ctx *gin.Context) {
	var err error
	var req dto.LoginRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	status, access, refresh, err := handler.userService.Login(&req)
	if err != nil || status != http.StatusOK {
		ctx.AbortWithStatus(status)
		return
	}

	ctx.SetCookie("refresh", refresh, jwt.RefreshMaxAge, "/api/auth/refresh", "", false, true)
	ctx.JSON(http.StatusOK, access)
}
