package repositories

import (
	"context"
	"time"

	"leetcode-spaced-repetition/models"

	"github.com/google/uuid"
)

type ProblemRepository interface {
	GetProblemSubmissions(c context.Context, problemIDs []int) ([]models.ProblemSubmissionWithDetails, error)
	GetSubmissionsByProblemID(c context.Context, problemID int) ([]models.ProblemSubmission, error)
	GetProblemByID(c context.Context, id int) (*models.Problem, error)
	GetProblemStatsByID(c context.Context, id int) (*models.ProblemSubmissionUserStats, error)
	GetAllProblemsPastReviewDate(c context.Context, limit uint) ([]models.Problem, error)
	GetProblems(c context.Context, topics, difficulties []string, page, limit int) ([]models.Problem, error)
	GetAllProblemTopics(c context.Context) ([]string, error)
	GetTopicsForProblem(c context.Context, id int) ([]string, error)
	SaveProblem(c context.Context, problem *models.Problem) error
	SaveProblemTopic(c context.Context, problemId int, topic string) error
	SaveProblemSubmission(c context.Context, problemID int, userID uuid.UUID, date time.Time, timeTaken *time.Duration, confidenceLevel models.ConfidenceLevel, language *models.SubmissionLanguage) (bool, error)
}
