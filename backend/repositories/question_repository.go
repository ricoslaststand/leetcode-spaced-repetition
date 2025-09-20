package repositories

import "leetcode-spaced-repetition/models"

type QuestionRepository interface {
	GetQuestionSubmissions() ([]models.QuestionSubmission, error)
	GetQuestionByID(id int) (*models.Question, error)
	GetQuestions() ([]models.Question, error)
	GetAllQuestionTags() ([]string, error)
	GetTagsForQuestion(ID int) ([]string, error)
	SaveQuestion(question models.Question) error
	SaveQuestionTag(questionId int, tag string) error
}
