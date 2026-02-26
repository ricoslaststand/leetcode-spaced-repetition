package controllers

import (
	"net/http"
	"regexp"
	"strconv"
	"time"

	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/services"

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
	TimeTaken       *int    `json:"timeTaken"`
	Language        *string `json:"language"`
	ConfidenceLevel string  `json:"confidenceLevel" binding:"required"`
	ProblemID       int     `json:"problemId" binding:"required,number" validate:"gte=1"`
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
		if err := v.RegisterValidation("date", validDate); err != nil {
			panic("failed to register date validator: " + err.Error())
		}
	}

	problemsController := ProblemsController{problemsService: *problemsService, logger: logger}

	r.GET("/dashboard", problemsController.GetDashboard)

	problemsGroup := r.Group("/problems")
	problemsGroup.GET("", problemsController.getProblems)
	problemsGroup.GET("topics", problemsController.GetAllProblemTopics)
	problemsGroup.GET("submissions", problemsController.getAllProblemSubmissions)
	problemsGroup.POST("submissions", problemsController.SaveProblemSubmission)
	problemsGroup.POST("submissions/import", problemsController.ImportSubmissions)
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
	difficulties, _ := ctx.GetQueryArray("difficulties")

	problems, err := c.problemsService.GetProblems(ctx, topics, difficulties, 1, 100)
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

// getProblemByID godoc
// @Summary      Get a problem by ID
// @Description  Returns a single problem including its topics
// @Tags         problems
// @Produce      json
// @Param        id   path  int  true  "Problem ID"
// @Success      200  {object}  models.ProblemDetail
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

	userID := uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")
	cardState, err := c.problemsService.GetProblemCardState(context, userID, problemID)
	if err != nil {
		c.logger.Error("failed to get card state for problem", zap.Int("problemID", problemID), zap.Error(err))
		context.JSON(500, gin.H{
			"error": "An internal server error has occurred.",
		})
		return
	}

	var stability *float64
	if cardState != nil {
		s := cardState.Stability
		stability = &s
	}

	context.JSON(200, models.ProblemDetail{Problem: *problem, Stability: stability})
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

	var timeTaken *time.Duration
	if problemSubmissionRequest.TimeTaken != nil {
		d := time.Duration(*problemSubmissionRequest.TimeTaken) * time.Second
		timeTaken = &d
	}

	var language *models.SubmissionLanguage
	if problemSubmissionRequest.Language != nil {
		l := models.SubmissionLanguage(*problemSubmissionRequest.Language)
		language = &l
	}

	if err := c.problemsService.SaveProblemSubmission(
		context,
		problemSubmissionRequest.ProblemID,
		uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11"),
		time.Now(),
		timeTaken,
		models.ConfidenceLevel(problemSubmissionRequest.ConfidenceLevel),
		language,
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

// ImportSubmissions godoc
// @Summary      Import submissions from Excel
// @Description  Accepts a multipart/form-data .xlsx file and bulk-imports problem submissions. Each row must have columns: "Problem #", "Date" (YYYY-MM-DD), "Confidence" (1-4), "Language" (python only).
// @Tags         submissions
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "Excel file (.xlsx)"
// @Success      200   {object}  models.ImportSubmissionsResult
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /problems/submissions/import [post]
func (c ProblemsController) ImportSubmissions(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "A file field named \"file\" is required.",
		})
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		c.logger.Error("failed to open uploaded file", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to read the uploaded file.",
		})
		return
	}
	defer func() { _ = f.Close() }()

	result, err := c.problemsService.ImportSubmissions(ctx, f)
	if err != nil {
		c.logger.Error("failed to import submissions", zap.Error(err))
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, result)
}

// GetDashboard godoc
// @Summary      Dashboard summary
// @Description  Returns problems due for review and problems with the lowest FSRS stability
// @Tags         dashboard
// @Produce      json
// @Param        limit  query     int  false  "Number of items per table (1–50, default 10)"
// @Success      200  {object}  models.DashboardData
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /dashboard [get]
func (c ProblemsController) GetDashboard(ctx *gin.Context) {
	userID := uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")

	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 50 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "limit must be between 1 and 50"})
		return
	}

	data, err := c.problemsService.GetDashboard(ctx, userID, limit)
	if err != nil {
		c.logger.Error("failed to get dashboard data", zap.Error(err))
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "An internal server error has occurred.",
		})
		return
	}

	ctx.JSON(http.StatusOK, data)
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
