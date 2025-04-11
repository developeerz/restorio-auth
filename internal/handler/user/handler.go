package user

import (
	"net/http"

	"github.com/developeerz/restorio-auth/internal/handler/user/dto"
	"github.com/developeerz/restorio-auth/internal/jwt"
	"github.com/gin-gonic/gin"
)

const jwtMaxAge = jwt.RefreshMaxAge

type Handler struct {
	service     Service
	refreshPath string
}

func NewHandler(service Service, refreshPath string) *Handler {
	return &Handler{service: service, refreshPath: refreshPath}
}

func (handler *Handler) SignUp(ctx *gin.Context) {
	var err error
	var req dto.SignUpRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	status, err := handler.service.SignUp(&req)
	if err != nil {
		ctx.AbortWithStatus(status)
		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *Handler) Verification(ctx *gin.Context) {
	var err error
	var req dto.VerificationRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	status, err := handler.service.Verify(&req)
	if err != nil {
		ctx.AbortWithStatus(status)
		return
	}

	ctx.Status(http.StatusOK)
}

func (handler *Handler) Login(ctx *gin.Context) {
	var err error
	var req dto.LoginRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	status, access, refresh, err := handler.service.Login(&req)
	if err != nil {
		ctx.AbortWithStatus(status)
		return
	}

	ctx.SetCookie("refresh", refresh, jwtMaxAge, handler.refreshPath, "", false, true)
	ctx.JSON(http.StatusOK, access)
}
