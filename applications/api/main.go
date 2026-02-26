// @title           Leetcode Buddy API
// @version         1.0
// @description     A spaced repetition API for tracking Leetcode practice sessions.
// @host            localhost:8000
// @BasePath        /
package main

import (
	"fmt"

	"leetcode-spaced-repetition/controllers"
	"leetcode-spaced-repetition/internal"
	"leetcode-spaced-repetition/repositories"
	"leetcode-spaced-repetition/services"

	_ "leetcode-spaced-repetition/docs"

	ginprometheus "github.com/zsais/go-gin-prometheus"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func main() {
	cfg, err := internal.GetConfig()
	if err != nil {
		panic("failed to load config")
	}

	logger := internal.NewLogger(cfg.AppEnv)
	defer logger.Sync() //nolint:errcheck // logger.Sync error is not actionable at shutdown

	db, err := internal.GetDBConnFromConfig(&cfg)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer func() { _ = db.Close() }()

	if err = db.Ping(); err != nil {
		logger.Error("failed to ping database", zap.Error(err))
	}

	problemsRepo := repositories.NewProblemPostgresRepository(db, logger)
	cardStateRepo := repositories.NewProblemCardStatePostgresRepository(db, logger)
	problemsService := services.NewProblemsService(problemsRepo, cardStateRepo, logger)

	router := gin.Default()
	router.Use(cors.Default()) // All origins allowed by default
	logger.Info("CORS is enabled")
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	p := ginprometheus.NewWithConfig(ginprometheus.Config{
		Subsystem: "gin",
	})
	p.Use(router)

	controllers.RegisterRoutes(router, problemsService, logger)

	// TODO: Turn this into a configurable port
	logger.Info("starting API server", zap.String("port", ":8000"))
	if err := router.Run(":8000"); err != nil {
		logger.Fatal("server failed", zap.Error(err))
	}
	fmt.Println()
}
