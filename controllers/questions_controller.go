package controllers

import (
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/services"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const dateRegexString = `^\d{4}-\d{2}-\d{2}T\d{2}`

type leetCodeQuestionRequest struct {
	ID int `uri:"id" binding:"required,number" validate:"gte=1"`
}

type getSubmissionsRequest struct {
	QuestionID []int `form:"questionId" collection_format:"csv" default:""`
}

type saveQuestionSubmissionRequest struct {
	QuestionID      int    `json:"questionId" binding:"required,number" validate:"gte=1"`
	TimeTaken       int    `json:"timeTaken" binding:"required,number" validate:"gte=0"`
	ConfidenceLevel string `json:"confidenceLevel" binding:"required"`
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
	logger           *zap.Logger
}

func RegisterRoutes(r *gin.Engine, questionsService *services.QuestionService, logger *zap.Logger) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("date", validDate)
	}

	questionsController := QuestionsController{questionsService: *questionsService, logger: logger}

	questionsGroup := r.Group("/questions")
	questionsGroup.GET("", questionsController.getQuestions)
	questionsGroup.GET("tags", questionsController.GetAllQuestionTags)
	questionsGroup.GET("submissions", questionsController.getAllQuestionSubmissions)
	questionsGroup.POST("submissions", questionsController.SaveQuestionSubmission)
	questionsGroup.GET(":id", questionsController.GetQuestionByID)

	individualQuestionsGroup := questionsGroup.Group(":id")
	individualQuestionsGroup.GET("submissions", questionsController.getQuestionSubmissions)
}

// getQuestions godoc
// @Summary      List questions
// @Description  Returns a paginated list of Leetcode questions, optionally filtered by tags
// @Tags         questions
// @Produce      json
// @Param        tags  query  []string  false  "Filter by tag names"  collectionFormat(csv)
// @Success      200   {object}  models.QuestionPage
// @Failure      500   {object}  map[string]string
// @Router       /questions [get]
func (c QuestionsController) getQuestions(ctx *gin.Context) {
	tags, _ := ctx.GetQueryArray("tags")

	questions, err := c.questionsService.GetQuestions(ctx, tags, 1, 100)
	if err != nil {
		c.logger.Error("failed to get questions", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "An error has occurred.",
		})
		return
	}

	resp := models.Pagaination[models.Question]{
		Data: questions,
	}

	if len(questions) == 0 {
		resp.Data = make([]models.Question, 0)
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetQuestionByID godoc
// @Summary      Get a question by ID
// @Description  Returns a single question including its tags
// @Tags         questions
// @Produce      json
// @Param        id   path  int  true  "Question ID"
// @Success      200  {object}  models.Question
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /questions/{id} [get]
func (c QuestionsController) GetQuestionByID(context *gin.Context) {
	var request leetCodeQuestionRequest
	if err := context.ShouldBindUri(&request); err != nil {
		context.JSON(400, gin.H{
			"error": "The id of the question must be a valid integer.",
		})
		return
	}

	questionID := request.ID

	question, err := c.questionsService.GetQuestionByID(context, questionID)
	if err != nil {
		c.logger.Error("failed to get question by ID", zap.Int("questionID", questionID), zap.Error(err))
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

	tags, err := c.questionsService.GetTagsForQuestion(context, questionID)
	if err != nil {
		c.logger.Error("failed to get tags for question", zap.Int("questionID", questionID), zap.Error(err))
		context.JSON(500, gin.H{
			"error": "An internal server error has occurred.",
		})
		return
	}

	question.Tags = tags

	context.JSON(200, *question)
}

// getQuestionSubmissions godoc
// @Summary      Get submissions for a question
// @Description  Returns all submissions for the question identified by id
// @Tags         questions
// @Produce      json
// @Param        id   path  int  true  "Question ID"
// @Success      200  {object}  models.QuestionSubmissionPage
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /questions/{id}/submissions [get]
func (c QuestionsController) getQuestionSubmissions(context *gin.Context) {
	var request leetCodeQuestionRequest
	if err := context.ShouldBindUri(&request); err != nil {
		context.JSON(400, gin.H{
			"error": "The id of the question must be a valid integer.",
		})
		return
	}

	result, err := c.questionsService.GetAllSubmissionsForQuestion(
		context,
		request.ID,
	)
	if err != nil {
		c.logger.Error("failed to get question submissions", zap.Int("questionID", request.ID), zap.Error(err))
		context.JSON(500, gin.H{
			"error": "An internal server error occurred.",
		})
		return
	}

	resp := models.Pagaination[models.QuestionSubmission]{
		Data: result,
	}

	if len(result) == 0 {
		resp.Data = make([]models.QuestionSubmission, 0)
	}

	context.JSON(200, resp)
}

// GetAllQuestionTags godoc
// @Summary      List all question tags
// @Description  Returns every tag that exists across all questions
// @Tags         questions
// @Produce      json
// @Success      200  {object}  map[string][]string
// @Failure      500  {object}  map[string]string
// @Router       /questions/tags [get]
func (c QuestionsController) GetAllQuestionTags(context *gin.Context) {
	tags, err := c.questionsService.GetAllQuestionTags(context)
	if err != nil {
		c.logger.Error("failed to get all question tags", zap.Error(err))
		context.JSON(500, gin.H{
			"error": "An internal server error has occurred.",
		})
		return
	}

	context.JSON(200, gin.H{
		"tags": tags,
	})
}

// SaveQuestionSubmission godoc
// @Summary      Record a new submission
// @Description  Saves a new attempt for a Leetcode question
// @Tags         submissions
// @Accept       json
// @Produce      json
// @Param        body  body  saveQuestionSubmissionRequest  true  "Submission details"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /questions/submissions [post]
func (c QuestionsController) SaveQuestionSubmission(context *gin.Context) {
	var questionSubmissionRequest saveQuestionSubmissionRequest
	if err := context.ShouldBindBodyWithJSON(&questionSubmissionRequest); err != nil {
		context.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := c.questionsService.SaveQuestionSubmission(
		context,
		questionSubmissionRequest.QuestionID,
		uuid.New(),
		time.Now(),
		time.Duration(questionSubmissionRequest.TimeTaken*int(time.Second)),
		models.ConfidenceLevel(questionSubmissionRequest.ConfidenceLevel),
	); err != nil {
		c.logger.Error("failed to save question submission", zap.Int("questionID", questionSubmissionRequest.QuestionID), zap.Error(err))
		context.JSON(500, gin.H{
			"error": "An internal error has occurred.",
		})
		return
	}

	context.JSON(201, gin.H{
		"message": "Successfully saved question submission",
	})
}

// getAllQuestionSubmissions godoc
// @Summary      List all submissions
// @Description  Returns submissions across all questions, optionally filtered by question IDs
// @Tags         submissions
// @Produce      json
// @Param        questionId  query  []int  false  "Filter by question IDs"  collectionFormat(csv)
// @Success      200         {object}  models.QuestionSubmissionWithDetailsPage
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /questions/submissions [get]
func (c QuestionsController) getAllQuestionSubmissions(context *gin.Context) {
	var request getSubmissionsRequest
	if err := context.ShouldBindQuery(&request); err != nil {
		context.JSON(400, gin.H{
			"error": "Invalid query parameters",
		})
	}

	questionIDs := request.QuestionID

	submissions, err := c.questionsService.GetQuestionSubmissions(context, questionIDs)
	if err != nil {
		c.logger.Error("failed to get all question submissions", zap.Error(err))
		context.JSON(500, gin.H{
			"error": "An internal server error has occurred.",
		})
	}

	resp := models.Pagaination[models.QuestionSubmissionWithDetails]{
		Data: submissions,
	}

	if len(submissions) == 0 {
		resp.Data = make([]models.QuestionSubmissionWithDetails, 0)
	}

	context.JSON(200, resp)
}
