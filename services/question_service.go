package services

import (
	"context"
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/repositories"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type QuestionService struct {
	questionRepo repositories.QuestionRepository
	logger       *zap.Logger
}

func NewQuestionsService(questionsRepo repositories.QuestionRepository, logger *zap.Logger) *QuestionService {
	return &QuestionService{
		questionRepo: questionsRepo,
		logger:       logger,
	}
}

func (s QuestionService) GetQuestions(ctx context.Context, tags []string, page int, limit int) ([]models.Question, error) {
	questions, err := s.questionRepo.GetQuestions(ctx, tags, page, limit)
	if err != nil {
		s.logger.Error("failed to get questions", zap.Error(err))
	}
	return questions, err
}

func (s QuestionService) GetQuestionByID(c context.Context, ID int) (*models.Question, error) {
	question, err := s.questionRepo.GetQuestionByID(c, ID)
	if err != nil {
		s.logger.Error("failed to get question by ID", zap.Int("questionID", ID), zap.Error(err))
	}
	return question, err
}

func (s QuestionService) GetQuestionSubmissions(c context.Context, questionIDs []int) ([]models.QuestionSubmissionWithDetails, error) {
	submissions, err := s.questionRepo.GetQuestionSubmissions(c, questionIDs)
	if err != nil {
		s.logger.Error("failed to get question submissions", zap.Error(err))
	}
	return submissions, err
}

func (s QuestionService) GetAllQuestionTags(c context.Context) ([]string, error) {
	tags, err := s.questionRepo.GetAllQuestionTags(c)
	if err != nil {
		s.logger.Error("failed to get all question tags", zap.Error(err))
	}
	return tags, err
}

func (s QuestionService) GetTagsForQuestion(c context.Context, ID int) ([]string, error) {
	tags, err := s.questionRepo.GetTagsForQuestion(c, ID)
	if err != nil {
		s.logger.Error("failed to get tags for question", zap.Int("questionID", ID), zap.Error(err))
	}
	return tags, err
}

func (s QuestionService) GetAllQuestionsPastReviewDate(c context.Context, limit uint) ([]models.Question, error) {
	questions, err := s.questionRepo.GetAllQuestionsPastReviewDate(c, limit)
	if err != nil {
		s.logger.Error("failed to get questions past review date", zap.Error(err))
	}
	return questions, err
}

func (s QuestionService) GetAllSubmissionsForQuestion(c context.Context, questionID int) ([]models.QuestionSubmission, error) {
	submissions, err := s.questionRepo.GetSubmissionsByQuestionID(c, questionID)
	if err != nil {
		s.logger.Error("failed to get submissions for question", zap.Int("questionID", questionID), zap.Error(err))
		return []models.QuestionSubmission{}, err
	}

	return submissions, nil
}

func (s QuestionService) SaveQuestionSubmission(
	c context.Context,
	questionID int,
	userID uuid.UUID,
	date time.Time,
	timeTaken time.Duration,
	confidenceLevel models.ConfidenceLevel,
) error {
	err := s.questionRepo.SaveQuestionSubmission(c, questionID, userID, date, timeTaken, confidenceLevel)
	if err != nil {
		s.logger.Error("failed to save question submission", zap.Int("questionID", questionID), zap.Error(err))
	}
	return err
}
