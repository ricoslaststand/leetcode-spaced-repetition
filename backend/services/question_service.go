package services

import (
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/repositories"
)

type QuestionService struct {
	questionRepo repositories.QuestionRepository
}

func NewQuestionsService(questionsRepo repositories.QuestionRepository) *QuestionService {
	return &QuestionService{
		questionRepo: questionsRepo,
	}
}

func (s QuestionService) GetQuestionByID(ID int) (*models.Question, error) {
	return s.questionRepo.GetQuestionByID(ID)
}

func (s QuestionService) GetQuestionSubmissions() ([]models.QuestionSubmission, error) {
	return s.questionRepo.GetQuestionSubmissions()
}

func (s QuestionService) GetAllQuestionTags() ([]string, error) {
	return s.questionRepo.GetAllQuestionTags()
}

func (s QuestionService) GetTagsForQuestion(ID int) ([]string, error) {
	return s.questionRepo.GetTagsForQuestion(ID)
}
