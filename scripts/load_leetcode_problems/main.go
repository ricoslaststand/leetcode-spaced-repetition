package main

import (
	"context"
	"encoding/json"
	"fmt"
	"leetcode-spaced-repetition/internal"
	config "leetcode-spaced-repetition/internal"
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/repositories"
	"os"
	"strconv"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Difficulty struct {
	Level int `json:"level"`
}

type Stat struct {
	QuestionID              int     `json:"question_id"`
	QuestionArticleLive     *bool   `json:"question__article__live"`
	QuestionArticleSlug     *string `json:"question__article__slug"`
	QuestionArticleHasVideo *bool   `json:"question__article__has_video_solution"`
	QuestionTitle           string  `json:"question__title"`
	QuestionTitleSlug       string  `json:"question__title_slug"`
	QuestionHide            bool    `json:"question__hide"`
	TotalACs                int     `json:"total_acs"`
	TotalSubmitted          int     `json:"total_submitted"`
	FrontendQuestionID      int     `json:"frontend_question_id"`
	IsNewQuestion           bool    `json:"is_new_question"`
}

type StatStatusPair struct {
	Stat       Stat       `json:"stat"`
	Status     *string    `json:"status"` // Could be null (use *string for nullable field)
	Difficulty Difficulty `json:"difficulty"`
	PaidOnly   bool       `json:"paid_only"`
	IsFavor    bool       `json:"is_favor"`
	Frequency  int        `json:"frequency"`
	Progress   int        `json:"progress"`
}

type ProblemData struct {
	UserName        string           `json:"user_name"`
	NumSolved       int              `json:"num_solved"`
	NumTotal        int              `json:"num_total"`
	AcEasy          int              `json:"ac_easy"`
	AcMedium        int              `json:"ac_medium"`
	AcHard          int              `json:"ac_hard"`
	StatStatusPairs []StatStatusPair `json:"stat_status_pairs"`
}

type Questions struct {
	Questions []struct {
		ID   string   `json:"problem_id"`
		Tags []string `json:"topics"`
	} `json:"questions"`
}

func main() {
	logger := internal.NewLogger()
	defer logger.Sync() //nolint:errcheck

	config, err := config.GetConfig()
	if err != nil {
		logger.Fatal("failed to load config", zap.Error(err))
	}

	db, err := internal.GetDBConnFromConfig(config)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	logger.Info("loading leetcode problems")

	filepath := "leetcode_problems.json"
	fileContent, err := os.ReadFile(filepath)
	if err != nil {
		logger.Fatal("failed to read problems file", zap.String("filepath", filepath), zap.Error(err))
	}

	var responseData ProblemData
	if err = json.Unmarshal(fileContent, &responseData); err != nil {
		logger.Fatal("failed to unmarshal problems JSON", zap.Error(err))
	}

	numOfQuestions := len(responseData.StatStatusPairs)
	logger.Info("downloaded problems", zap.Int("count", numOfQuestions))

	questionRepo := repositories.NewQuestionPostgresRepository(db, logger)

	var questions []models.Question
	for i := 0; i < numOfQuestions; i++ {
		currQuestion := responseData.StatStatusPairs[i]
		difficultyMap := map[int]models.QuestionDifficulty{
			1: models.EasyDifficulty,
			2: models.MildDifficulty,
			3: models.HardDifficulty,
		}
		questionDifficulty, ok := difficultyMap[currQuestion.Difficulty.Level]
		if !ok {
			logger.Error("unrecognized difficulty level", zap.Int("level", currQuestion.Difficulty.Level))
			return
		}

		questions = append(questions, models.Question{
			ID:          currQuestion.Stat.QuestionID,
			Title:       currQuestion.Stat.QuestionTitle,
			Description: "",
			Slug:        currQuestion.Stat.QuestionTitleSlug,
			Difficulty:  questionDifficulty,
		})
	}

	logger.Info("saving questions to the database", zap.Int("count", len(questions)))

	c := context.Background()

	for i := 0; i < len(questions); i++ {
		if err = questionRepo.SaveQuestion(c, &questions[i]); err != nil {
			logger.Fatal("failed to save question", zap.Int("questionID", questions[i].ID), zap.Error(err))
		}
	}

	// Adding question tags
	tagsFileContent, err := os.ReadFile("merged_problems.json")
	if err != nil {
		logger.Fatal("failed to read merged_problems.json", zap.Error(err))
	}

	var mergedQuestions Questions
	if err = json.Unmarshal(tagsFileContent, &mergedQuestions); err != nil {
		logger.Fatal("failed to unmarshal tags JSON", zap.Error(err))
	}

	for i := 0; i < len(mergedQuestions.Questions); i++ {
		question := mergedQuestions.Questions[i]
		questionId, err := strconv.Atoi(question.ID)
		if err != nil {
			logger.Error("invalid question ID", zap.String("id", question.ID), zap.Error(err))
			continue
		}

		for j := 0; j < len(question.Tags); j++ {
			logger.Debug("processing tag", zap.Int("questionID", questionId), zap.String("tag", question.Tags[j]))
			if err = questionRepo.SaveQuestionTag(c, questionId, question.Tags[j]); err != nil {
				logger.Fatal("failed to save question tag", zap.Int("questionID", questionId), zap.String("tag", question.Tags[j]), zap.Error(err))
			}
		}
	}

	logger.Info("load_leetcode_problems completed successfully")
	fmt.Println()
}
