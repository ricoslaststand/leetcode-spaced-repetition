package services

import (
	"context"
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/repositories"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type ProblemService struct {
	problemRepo repositories.ProblemRepository
	logger      *zap.Logger
}

func NewProblemsService(problemsRepo repositories.ProblemRepository, logger *zap.Logger) *ProblemService {
	return &ProblemService{
		problemRepo: problemsRepo,
		logger:      logger,
	}
}

func (s ProblemService) GetProblems(ctx context.Context, topics []string, page int, limit int) ([]models.Problem, error) {
	problems, err := s.problemRepo.GetProblems(ctx, topics, page, limit)
	if err != nil {
		s.logger.Error("failed to get problems", zap.Error(err))
	}
	return problems, err
}

func (s ProblemService) GetProblemByID(c context.Context, ID int) (*models.Problem, error) {
	problem, err := s.problemRepo.GetProblemByID(c, ID)
	if err != nil {
		s.logger.Error("failed to get problem by ID", zap.Int("problemID", ID), zap.Error(err))
	}
	return problem, err
}

func (s ProblemService) GetProblemSubmissions(c context.Context, problemIDs []int) ([]models.ProblemSubmissionWithDetails, error) {
	submissions, err := s.problemRepo.GetProblemSubmissions(c, problemIDs)
	if err != nil {
		s.logger.Error("failed to get problem submissions", zap.Error(err))
	}
	return submissions, err
}

func (s ProblemService) GetAllProblemTopics(c context.Context) ([]string, error) {
	topics, err := s.problemRepo.GetAllProblemTopics(c)
	if err != nil {
		s.logger.Error("failed to get all problem topics", zap.Error(err))
	}
	return topics, err
}

func (s ProblemService) GetTopicsForProblem(c context.Context, ID int) ([]string, error) {
	topics, err := s.problemRepo.GetTopicsForProblem(c, ID)
	if err != nil {
		s.logger.Error("failed to get topics for problem", zap.Int("problemID", ID), zap.Error(err))
	}
	return topics, err
}

func (s ProblemService) GetAllProblemsPastReviewDate(c context.Context, limit uint) ([]models.Problem, error) {
	problems, err := s.problemRepo.GetAllProblemsPastReviewDate(c, limit)
	if err != nil {
		s.logger.Error("failed to get problems past review date", zap.Error(err))
	}
	return problems, err
}

func (s ProblemService) GetAllSubmissionsForProblem(c context.Context, problemID int) ([]models.ProblemSubmission, error) {
	submissions, err := s.problemRepo.GetSubmissionsByProblemID(c, problemID)
	if err != nil {
		s.logger.Error("failed to get submissions for problem", zap.Int("problemID", problemID), zap.Error(err))
		return []models.ProblemSubmission{}, err
	}

	return submissions, nil
}

func (s ProblemService) SaveProblemSubmission(
	c context.Context,
	problemID int,
	userID uuid.UUID,
	date time.Time,
	timeTaken time.Duration,
	confidenceLevel models.ConfidenceLevel,
) error {
	err := s.problemRepo.SaveProblemSubmission(c, problemID, userID, date, timeTaken, confidenceLevel)
	if err != nil {
		s.logger.Error("failed to save problem submission", zap.Int("problemID", problemID), zap.Error(err))
	}
	return err
}
