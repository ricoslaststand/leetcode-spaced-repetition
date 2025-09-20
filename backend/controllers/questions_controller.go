package controllers

import (
	"fmt"
	"leetcode-spaced-repetition/services"
	"regexp"
	"strconv"

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

type QuestionsController struct {
	questionsService services.QuestionService
}

func RegisterRoutes(r *gin.Engine, questionsService *services.QuestionService) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("date", validDate)
	}

	questionsController := QuestionsController{questionsService: *questionsService}

	questionsGroup := r.Group("/questions")
	questionsGroup.GET("tags", questionsController.GetAllQuestionTags)
	questionsGroup.GET(":id", questionsController.GetQuestionByID)

	individualQuestionsGroup := questionsGroup.Group(":id")
	individualQuestionsGroup.GET("submissions", getQuestionSubmissionsByID)
}

func (c QuestionsController) GetQuestionByID(context *gin.Context) {
	var request leetCodeQuestionRequest
	if err := context.ShouldBindUri(&request); err != nil {
		context.JSON(400, gin.H{
			"error": "The id of the question must be a valid integer.",
		})
		return
	}

	intId, err := strconv.Atoi(request.ID)
	if err != nil {
		context.JSON(400, gin.H{
			"error": "The id of the question must be a valid integer.",
		})
		return
	}

	question, err := c.questionsService.GetQuestionByID(intId)
	if err != nil {
		context.JSON(500, gin.H{
			"error": "An internal server has occurred.",
		})
		return
	}
	if question == nil {
		context.JSON(404, gin.H{
			"message": "No question is associated with this code.",
		})
		return
	}

	tags, err := c.questionsService.GetTagsForQuestion(intId)
	if err != nil {
		context.JSON(500, gin.H{
			"error": "An internal server error has occurred.",
		})
		return
	}

	question.Tags = tags

	fmt.Printf("The question is %+v\n", *question)

	context.JSON(200, *question)
}

func getQuestionSubmissionsByID(c *gin.Context) {
	var request questionSubmissionRequest

	if err := c.ShouldBindUri(&request); err != nil {

	}

	c.JSON(200, gin.H{
		"data": []string{},
	})
}

func (c QuestionsController) GetAllQuestionTags(context *gin.Context) {
	tags, err := c.questionsService.GetAllQuestionTags()
	if err != nil {
		context.JSON(500, gin.H{
			"error": "An internal server error has occurred.",
		})
	}

	context.JSON(200, gin.H{
		"tags": tags,
	})
}
