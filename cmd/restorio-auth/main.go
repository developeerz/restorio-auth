package main

import (
	"fmt"

	"github.com/developeerz/restorio-auth/config"
	"github.com/developeerz/restorio-auth/internal/database"
	auth_handler "github.com/developeerz/restorio-auth/internal/handler/auth"
	user_handler "github.com/developeerz/restorio-auth/internal/handler/user"
	"github.com/developeerz/restorio-auth/internal/middleware"
	"github.com/developeerz/restorio-auth/internal/repository/postgres"
	"github.com/developeerz/restorio-auth/internal/routers"
	auth_service "github.com/developeerz/restorio-auth/internal/service/auth"
	user_service "github.com/developeerz/restorio-auth/internal/service/user"
	"github.com/developeerz/restorio-auth/logger"
	"github.com/developeerz/restorio-auth/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func main() {
	err := logger.InitLogger()
	if err != nil {
		log.Fatal().AnErr("error", err)
	}

	config.LoadConfig()

	pgdb, err := database.PostgresConnect()
	if err != nil {
		log.Fatal().AnErr("error", err)
	}

	rdb, err := database.RedisConnect()
	if err != nil {
		log.Fatal().AnErr("error", err)
	}

	userRepo := postgres.NewUserRepository(pgdb)
	userCache := redis.NewUserCache(rdb)

	userService := user_service.NewService(userRepo, userCache)
	userHandler := user_handler.NewHandler(userService, routers.AuthGroupFullRefreshPath)

	authService := auth_service.NewService(userRepo)
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
