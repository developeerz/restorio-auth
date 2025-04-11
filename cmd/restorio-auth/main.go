package main

import (
	"log"

	"github.com/developeerz/restorio-auth/config"
	"github.com/developeerz/restorio-auth/internal/database"
	auth_handler "github.com/developeerz/restorio-auth/internal/handler/auth"
	user_handler "github.com/developeerz/restorio-auth/internal/handler/user"
	"github.com/developeerz/restorio-auth/internal/repository"
	"github.com/developeerz/restorio-auth/internal/routers"
	auth_service "github.com/developeerz/restorio-auth/internal/service/auth"
	user_service "github.com/developeerz/restorio-auth/internal/service/user"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadConfig()

	db := database.Connect()
	repository := repository.NewRepository(db)

	userService := user_service.NewService(repository)
	userHandler := user_handler.NewHandler(userService, routers.AuthGroupFullRefreshPath)

	authService := auth_service.NewService(repository)
	authHandler := auth_handler.NewHandler(authService, routers.AuthGroupFullRefreshPath)

	router := gin.Default()
	routers.NewUserRouter(router, userHandler)
	routers.NewAuthRouter(router, authHandler)

	err := router.Run(":8081")
	if err != nil {
		log.Fatalf("start server error: %v", err)
	}
}
