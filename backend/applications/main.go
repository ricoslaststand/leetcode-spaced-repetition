package main

import (
	"leetcode-spaced-repetition/controllers"
	config "leetcode-spaced-repetition/internal"

	ginprometheus "github.com/zsais/go-gin-prometheus"

	"github.com/gin-gonic/gin"
)

func main() {
	_, err := config.GetConfig()
	if err != nil {
		panic(err)
	}

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

	controllers.RegisterRoutes(router)

	router.Run()
}
