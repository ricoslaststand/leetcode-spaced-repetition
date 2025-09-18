package controllers

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

const dateRegexString = `^\d{4}-\d{2}-\d{2}T\d{2}`

type leetCodeQuestionRequest struct {
	ID string `uri:"id" binding:"required,number"`
}

type questionSubmissionRequest struct {
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
}

var validDate validator.Func = func(fl validator.FieldLevel) bool {
	date, ok := fl.Field().Interface().(string)

	if !ok {
		return false
	}
	match, _ := regexp.MatchString(dateRegexString, date)
	return match
}

func RegisterRoutes(r *gin.Engine) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("date", validDate)
	}

	questionsGroup := r.Group("/questions")
	questionsGroup.GET(":id", getQuestionByID)

	individualQuestionsGroup := questionsGroup.Group(":id")
	individualQuestionsGroup.GET("submissions", getQuestionSubmissionsByID)
}

func getQuestionByID(c *gin.Context) {
	var request leetCodeQuestionRequest
	if err := c.ShouldBindUri(&request); err != nil {
		c.JSON(400, gin.H{
			"error": "The id of the question must be a valid integer.",
		})
	}

	c.JSON(200, gin.H{
		"message": "We will get back to you with an answer.",
	})
}

func getQuestionSubmissionsByID(c *gin.Context) {
	var request questionSubmissionRequest

	if err := c.ShouldBindUri(&request); err != nil {

	}

	c.JSON(200, gin.H{
		"data": []string{},
	})
}
