package main

import (
	"leetcode-spaced-repetition/controllers"
	"leetcode-spaced-repetition/internal"
	config "leetcode-spaced-repetition/internal"
	"leetcode-spaced-repetition/repositories"
	"leetcode-spaced-repetition/services"
	"log"

	ginprometheus "github.com/zsais/go-gin-prometheus"

	"github.com/gin-gonic/gin"
)

func main() {
	config, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

	db, err := internal.GetDBConnFromConfig(config)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	questionsRepo := repositories.NewQuestionPostgresRepository(db)
	questionsService := services.NewQuestionsService(questionsRepo)

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	p := ginprometheus.NewWithConfig(ginprometheus.Config{
		Subsystem: "gin",
	})
	p.Use(router)

	controllers.RegisterRoutes(router, questionsService)

	router.Run(":8000")
}
