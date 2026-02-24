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

type leetCodeProblemRequest struct {
	ID int `uri:"id" binding:"required,number" validate:"gte=1"`
}

type getProblemSubmissionsRequest struct {
	ProblemID []int `form:"problemId" collection_format:"csv" default:""`
}

type saveProblemSubmissionRequest struct {
	ProblemID       int    `json:"problemId" binding:"required,number" validate:"gte=1"`
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

type ProblemsController struct {
	problemsService services.ProblemService
	logger          *zap.Logger
}

func RegisterRoutes(r *gin.Engine, problemsService *services.ProblemService, logger *zap.Logger) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("date", validDate)
	}

	problemsController := ProblemsController{problemsService: *problemsService, logger: logger}

	problemsGroup := r.Group("/problems")
	problemsGroup.GET("", problemsController.getProblems)
	problemsGroup.GET("topics", problemsController.GetAllProblemTopics)
	problemsGroup.GET("submissions", problemsController.getAllProblemSubmissions)
	problemsGroup.POST("submissions", problemsController.SaveProblemSubmission)
	problemsGroup.GET(":id", problemsController.GetProblemByID)

	individualProblemsGroup := problemsGroup.Group(":id")
	individualProblemsGroup.GET("submissions", problemsController.getProblemSubmissions)
}

// getProblems godoc
// @Summary      List problems
// @Description  Returns a paginated list of Leetcode problems, optionally filtered by topics
// @Tags         problems
// @Produce      json
// @Param        topics  query  []string  false  "Filter by topic names"  collectionFormat(csv)
// @Success      200   {object}  models.ProblemPage
// @Failure      500   {object}  map[string]string
// @Router       /problems [get]
func (c ProblemsController) getProblems(ctx *gin.Context) {
	topics, _ := ctx.GetQueryArray("topics")

	problems, err := c.problemsService.GetProblems(ctx, topics, 1, 100)
	if err != nil {
		c.logger.Error("failed to get problems", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "An error has occurred.",
		})
		return
	}

	resp := models.Pagaination[models.Problem]{
		Data: problems,
	}

	if len(problems) == 0 {
		resp.Data = make([]models.Problem, 0)
	}

	ctx.JSON(http.StatusOK, resp)
}

// GetProblemByID godoc
// @Summary      Get a problem by ID
// @Description  Returns a single problem including its topics
// @Tags         problems
// @Produce      json
// @Param        id   path  int  true  "Problem ID"
// @Success      200  {object}  models.Problem
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /problems/{id} [get]
func (c ProblemsController) GetProblemByID(context *gin.Context) {
	var request leetCodeProblemRequest
	if err := context.ShouldBindUri(&request); err != nil {
		context.JSON(400, gin.H{
			"error": "The id of the problem must be a valid integer.",
		})
		return
	}

	problemID := request.ID

	problem, err := c.problemsService.GetProblemByID(context, problemID)
	if err != nil {
		c.logger.Error("failed to get problem by ID", zap.Int("problemID", problemID), zap.Error(err))
		context.JSON(500, gin.H{
			"error": "An internal server has occurred.",
		})
		return
	}
	if problem == nil {
		context.JSON(404, gin.H{
			"message": "No problem is associated with this code.",
		})
		return
	}

	topics, err := c.problemsService.GetTopicsForProblem(context, problemID)
	if err != nil {
		c.logger.Error("failed to get topics for problem", zap.Int("problemID", problemID), zap.Error(err))
		context.JSON(500, gin.H{
			"error": "An internal server error has occurred.",
		})
		return
	}

	problem.Topics = topics

	context.JSON(200, *problem)
}

// getProblemSubmissions godoc
// @Summary      Get submissions for a problem
// @Description  Returns all submissions for the problem identified by id
// @Tags         problems
// @Produce      json
// @Param        id   path  int  true  "Problem ID"
// @Success      200  {object}  models.ProblemSubmissionPage
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /problems/{id}/submissions [get]
func (c ProblemsController) getProblemSubmissions(context *gin.Context) {
	var request leetCodeProblemRequest
	if err := context.ShouldBindUri(&request); err != nil {
		context.JSON(400, gin.H{
			"error": "The id of the problem must be a valid integer.",
		})
		return
	}

	result, err := c.problemsService.GetAllSubmissionsForProblem(
		context,
		request.ID,
	)
	if err != nil {
		c.logger.Error("failed to get problem submissions", zap.Int("problemID", request.ID), zap.Error(err))
		context.JSON(500, gin.H{
			"error": "An internal server error occurred.",
		})
		return
	}

	resp := models.Pagaination[models.ProblemSubmission]{
		Data: result,
	}

	if len(result) == 0 {
		resp.Data = make([]models.ProblemSubmission, 0)
	}

	context.JSON(200, resp)
}

// GetAllProblemTopics godoc
// @Summary      List all problem topics
// @Description  Returns every topic that exists across all problems
// @Tags         problems
// @Produce      json
// @Success      200  {object}  map[string][]string
// @Failure      500  {object}  map[string]string
// @Router       /problems/topics [get]
func (c ProblemsController) GetAllProblemTopics(context *gin.Context) {
	topics, err := c.problemsService.GetAllProblemTopics(context)
	if err != nil {
		c.logger.Error("failed to get all problem topics", zap.Error(err))
		context.JSON(500, gin.H{
			"error": "An internal server error has occurred.",
		})
		return
	}

	context.JSON(200, gin.H{
		"topics": topics,
	})
}

// SaveProblemSubmission godoc
// @Summary      Record a new submission
// @Description  Saves a new attempt for a Leetcode problem
// @Tags         submissions
// @Accept       json
// @Produce      json
// @Param        body  body  saveProblemSubmissionRequest  true  "Submission details"
// @Success      201   {object}  map[string]string
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /problems/submissions [post]
func (c ProblemsController) SaveProblemSubmission(context *gin.Context) {
	var problemSubmissionRequest saveProblemSubmissionRequest
	if err := context.ShouldBindBodyWithJSON(&problemSubmissionRequest); err != nil {
		context.JSON(400, gin.H{
			"error": "Invalid request body",
		})
		return
	}

	if err := c.problemsService.SaveProblemSubmission(
		context,
		problemSubmissionRequest.ProblemID,
		uuid.New(),
		time.Now(),
		time.Duration(problemSubmissionRequest.TimeTaken*int(time.Second)),
		models.ConfidenceLevel(problemSubmissionRequest.ConfidenceLevel),
	); err != nil {
		c.logger.Error("failed to save problem submission", zap.Int("problemID", problemSubmissionRequest.ProblemID), zap.Error(err))
		context.JSON(500, gin.H{
			"error": "An internal error has occurred.",
		})
		return
	}

	context.JSON(201, gin.H{
		"message": "Successfully saved problem submission",
	})
}

// getAllProblemSubmissions godoc
// @Summary      List all submissions
// @Description  Returns submissions across all problems, optionally filtered by problem IDs
// @Tags         submissions
// @Produce      json
// @Param        problemId  query  []int  false  "Filter by problem IDs"  collectionFormat(csv)
// @Success      200         {object}  models.ProblemSubmissionWithDetailsPage
// @Failure      400         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /problems/submissions [get]
func (c ProblemsController) getAllProblemSubmissions(context *gin.Context) {
	var request getProblemSubmissionsRequest
	if err := context.ShouldBindQuery(&request); err != nil {
		context.JSON(400, gin.H{
			"error": "Invalid query parameters",
		})
	}

	problemIDs := request.ProblemID

	submissions, err := c.problemsService.GetProblemSubmissions(context, problemIDs)
	if err != nil {
		c.logger.Error("failed to get all problem submissions", zap.Error(err))
		context.JSON(500, gin.H{
			"error": "An internal server error has occurred.",
		})
	}

	resp := models.Pagaination[models.ProblemSubmissionWithDetails]{
		Data: submissions,
	}

	if len(submissions) == 0 {
		resp.Data = make([]models.ProblemSubmissionWithDetails, 0)
	}

	context.JSON(200, resp)
}
