package handler

import (
	"net/http"
	"strings"

	"github.com/developeerz/restorio-auth/internal/dto"
	"github.com/developeerz/restorio-auth/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: *userService}
}

func (handler *UserHandler) SignUp(ctx *gin.Context) {
	var err error
	var req dto.SignUpRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Error{Message: "Cannot parse json"})
		return
	}

	req.Telegram, _ = strings.CutPrefix(req.Telegram, "@")

	status, msg, err := handler.userService.SignUp(&req)
	if err != nil || status != http.StatusOK {
		ctx.JSON(status, msg)
		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *UserHandler) Verification(ctx *gin.Context) {
	var err error
	var req dto.VerificationRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Error{Message: "Cannot parse json"})
		return
	}

	req.Telegram, _ = strings.CutPrefix(req.Telegram, "@")

	status, msg, err := handler.userService.Verify(&req)
	if err != nil {
		ctx.JSON(status, msg)
		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *UserHandler) SignIn(ctx *gin.Context) {
	var err error
	var req dto.SignInRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, dto.Error{Message: "Cannot parse json"})
		return
	}

	req.Telegram, _ = strings.CutPrefix(req.Telegram, "@")

	status, jwt, msg, err := handler.userService.SignIn(&req)
	if err != nil || status != http.StatusOK {
		ctx.JSON(status, msg)
		return
	}

	ctx.JSON(http.StatusOK, jwt)
}
