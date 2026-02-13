package repositories

import (
	"context"
	"leetcode-spaced-repetition/models"
	"time"

	"github.com/google/uuid"
)

type QuestionRepository interface {
	GetQuestionSubmissions(c context.Context, questionIDs []int) ([]models.QuestionSubmissionWithDetails, error)
	GetSubmissionsByQuestionID(c context.Context, questionID int) ([]models.QuestionSubmission, error)
	GetQuestionByID(c context.Context, id int) (*models.Question, error)
	GetQuestionStatsByID(c context.Context, id int) (*models.QuestionSubmissionUserStats, error)
	GetAllQuestionsPastReviewDate(c context.Context, limit uint) ([]models.Question, error)
	GetQuestions(c context.Context, tags []string, page, limit int) ([]models.Question, error)
	GetAllQuestionTags(c context.Context) ([]string, error)
	GetTagsForQuestion(c context.Context, ID int) ([]string, error)
	SaveQuestion(c context.Context, question *models.Question) error
	SaveQuestionTag(c context.Context, questionId int, tag string) error
	SaveQuestionSubmission(c context.Context, questionID int, userID uuid.UUID, date time.Time, timeTaken time.Duration, confidenceLevel models.ConfidenceLevel) error
}
