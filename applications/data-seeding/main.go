package main

import (
	"context"
	"encoding/json"
	"fmt"
	"leetcode-spaced-repetition/internal"
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/repositories"
	"os"
	"strconv"
	"strings"

	"github.com/mgechev/revive/config"
	"go.uber.org/zap"
)

type LeetcodeProblem struct {
	Name       string
	Difficulty models.QuestionDifficulty
	Slug       string
	Acceptance float64
	Frequency  float64
	TimeTag    string
}

type questionsResponse struct {
	Questions []struct {
		Title       string   `json:"title"`
		ProblemID   string   `json:"problem_id"`
		FrontendID  string   `json:"frontend_id"`
		Difficulty  string   `json:"difficulty"`
		ProblemSlug string   `json:"problem_slug"`
		Topics      []string `json:"topics"`
		Description string   `json:"description"`
		Constraints []string `json:"constraints"`
		FollowUps   []string `json:"follow_ups"`
		Hints       []string `json:"hints"`
		Solution    string   `json:"solution"`
	} `json:"questions"`
}

type TimeTag string

const (
	LastThirtyDays    = "last_thirty_days"
	LastThreeMonths   = "last_three_months"
	LastSixMonths     = "last_six_months"
	MoreThanSixMonths = "more_than_six_months"
)

const (
	DifficultyIdx int = iota
	NameIdx
	FrequencyIdx
	AcceptanceIdx
	LinkIdx
)

func convertFileNameToTimeTag(filename string) (TimeTag, error) {
	lCase := strings.ToLower(filename)

	if strings.Contains(lCase, "more than six months") {
		return MoreThanSixMonths, nil
	} else if strings.Contains(lCase, "six months") {
		return LastSixMonths, nil
	} else if strings.Contains(lCase, "three months") {
		return LastThreeMonths, nil
	} else if strings.Contains(lCase, "thirty days") {
		return LastThirtyDays, nil
	}
	return LastThirtyDays, fmt.Errorf("'%s' is not a valid time tag", filename)
}

func convertStringToDifficulty(diffStr string) (models.QuestionDifficulty, error) {
	switch strings.ToLower(strings.Trim(diffStr, " `")) {
	case "easy":
		return models.EasyDifficulty, nil
	case "medium", "mild":
		return models.MildDifficulty, nil
	case "hard":
		return models.HardDifficulty, nil
	default:
		return models.EasyDifficulty, fmt.Errorf("'%s' is not a valid difficulty", diffStr)
	}
}

func getSlugFromLink(link string) string {
	parts := strings.Split(link, "/")

	return parts[len(parts)-1]
}

func main() {
	config, err := config.GetConfig()
	if err != nil {
		panic("failed to load config")
	}

	logger := internal.NewLogger(config.AppEnv)
	defer logger.Sync() //nolint:errcheck

	db, err := internal.GetDBConnFromConfig(config)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	questionsRepo := repositories.NewQuestionPostgresRepository(db, logger)

	entries, err := os.ReadDir("./")
	if err != nil {
		logger.Fatal("failed to read directory", zap.Error(err))
	}

	for _, e := range entries {
		logger.Info("found file", zap.String("name", e.Name()))
	}

	tagsFileContent, err := os.ReadFile("merged_problems.json")
	if err != nil {
		logger.Error("failed to read merged_problems.json", zap.Error(err))
		return
	}

	var mergedQuestions questionsResponse
	if err = json.Unmarshal(tagsFileContent, &mergedQuestions); err != nil {
		logger.Error("failed to unmarshal questions JSON", zap.Error(err))
		return
	}

	for _, question := range mergedQuestions.Questions {
		questionDifficulty, err := convertStringToDifficulty(question.Difficulty)
		if err != nil {
			logger.Error("invalid difficulty", zap.String("difficulty", question.Difficulty), zap.Error(err))
			return
		}

		intQuestionID, err := strconv.Atoi(question.ProblemID)
		if err != nil {
			logger.Error("invalid problem ID", zap.String("problemID", question.ProblemID), zap.Error(err))
			return
		}

		if err = questionsRepo.SaveQuestion(context.Background(), &models.Question{
			ID:          intQuestionID,
			Title:       question.Title,
			Description: question.Description,
			Slug:        question.ProblemSlug,
			Difficulty:  questionDifficulty,
			Tags:        question.Topics,
		}); err != nil {
			logger.Error("failed to save question", zap.Int("questionID", intQuestionID), zap.Error(err))
			return
		}

		for _, tag := range question.Topics {
			if err = questionsRepo.SaveQuestionTag(context.Background(), intQuestionID, tag); err != nil {
				logger.Error("failed to save question tag", zap.Int("questionID", intQuestionID), zap.String("tag", tag), zap.Error(err))
				return
			}
		}
	}

	logger.Info("data seeding completed successfully")
}
