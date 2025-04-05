package main

import (
	"github.com/developeerz/restorio-auth/config"
	"github.com/developeerz/restorio-auth/internal/database"
	"github.com/developeerz/restorio-auth/internal/handler"
	"github.com/developeerz/restorio-auth/internal/middleware"
	"github.com/developeerz/restorio-auth/internal/repository"
	"github.com/developeerz/restorio-auth/internal/routers"
	"github.com/developeerz/restorio-auth/internal/service/auth"
	"github.com/developeerz/restorio-auth/internal/service/user"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database.Connect()

	repository := repository.NewRepository(database.DB)

	userService := user.NewUserService(repository)
	userHandler := handler.NewUserHandler(userService)

	authService := auth.NewAuthService(repository)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.Default()

	middleware.ConfigureCORS(router)

	routers.NewUserRouter(router, userHandler)
	routers.NewAuthRouter(router, authHandler)

	router.Run(":8081")
}
