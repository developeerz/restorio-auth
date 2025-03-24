package main

import (
	"github.com/developeerz/restorio-auth/config"
	"github.com/developeerz/restorio-auth/internal/database"
	"github.com/developeerz/restorio-auth/internal/handler"
	"github.com/developeerz/restorio-auth/internal/repository"
	"github.com/developeerz/restorio-auth/internal/routers"
	"github.com/developeerz/restorio-auth/internal/service"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()
	database.Connect()

	userRepository := repository.NewUserRipository(database.DB)

	userService := service.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)

	authService := service.NewAuthService(userRepository)
	authHandler := handler.NewAuthHandler(authService)

	router := gin.Default()
	routers.NewUserRouter(router, userHandler)
	routers.NewAuthRouter(router, authHandler)

	router.Run(":8081")
}
