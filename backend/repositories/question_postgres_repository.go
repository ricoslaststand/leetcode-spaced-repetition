package repositories

import (
	"database/sql"
	"leetcode-spaced-repetition/models"
)

type QuestionPostgresRepository struct {
	db *sql.DB
}

func NewQuestionPostgresRepository(db *sql.DB) *QuestionPostgresRepository {
	return &QuestionPostgresRepository{
		db: db,
	}
}

func (r QuestionPostgresRepository) GetQuestionByID(ID int) *models.Question {
	return &models.Question{}
}

func (r QuestionPostgresRepository) GetQuestionSubmissions() []models.QuestionSubmission {
	var submissions []models.QuestionSubmission

	return submissions
}

func (r QuestionPostgresRepository) SaveQuestion(q models.Question) error {
	return nil
}
