package main

import (
	"fmt"

	"github.com/developeerz/restorio-auth/config"
	"github.com/developeerz/restorio-auth/internal/database"
	auth_handler "github.com/developeerz/restorio-auth/internal/handler/auth"
	user_handler "github.com/developeerz/restorio-auth/internal/handler/user"
	"github.com/developeerz/restorio-auth/internal/middleware"
	"github.com/developeerz/restorio-auth/internal/repository"
	"github.com/developeerz/restorio-auth/internal/routers"
	auth_service "github.com/developeerz/restorio-auth/internal/service/auth"
	user_service "github.com/developeerz/restorio-auth/internal/service/user"
	"github.com/developeerz/restorio-auth/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	err := logger.InitLogger()
	if err != nil {
		log.Fatal().AnErr("error", fmt.Errorf("logger init: %w", err))
	}

	config.LoadConfig()

	db, err := database.Connect()
	if err != nil {
		log.Fatal().AnErr("error", fmt.Errorf("database connect: %w", err))
	}

	repository := repository.NewRepository(db)

	userService := user_service.NewService(repository)
	userHandler := user_handler.NewHandler(userService, routers.AuthGroupFullRefreshPath)

	authService := auth_service.NewService(repository)
	authHandler := auth_handler.NewHandler(authService, routers.AuthGroupFullRefreshPath)

	router := gin.Default()
	router.Use(middleware.Logging)

	routers.NewUserRouter(router, userHandler)
	routers.NewAuthRouter(router, authHandler)

	err = router.Run(":8081")
	if err != nil {
		log.Fatal().AnErr("error", fmt.Errorf("server run: %w", err))
	}
}
