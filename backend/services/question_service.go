package services

import (
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/repositories"
)

type QuestionService struct {
	questionRepo repositories.QuestionPostgresRepository
}

func (s QuestionService) GetQuestionByID(ID int) *models.Question {
	return s.questionRepo.GetQuestionByID(ID)
}

func (s QuestionService) GetQuestionSubmissions() []models.QuestionSubmission {

}
