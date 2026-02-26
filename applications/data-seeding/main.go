package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"leetcode-spaced-repetition/internal"
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/repositories"

	"go.uber.org/zap"
)

type LeetcodeProblem struct {
	Name       string
	Difficulty models.ProblemDifficulty
	Slug       string
	TimeTag    string
	Acceptance float64
	Frequency  float64
}

type problemsResponse struct {
	Problems []struct {
		Title       string   `json:"title"`
		ProblemID   string   `json:"problem_id"`
		FrontendID  string   `json:"frontend_id"`
		Difficulty  string   `json:"difficulty"`
		ProblemSlug string   `json:"problem_slug"`
		Description string   `json:"description"`
		Solution    string   `json:"solution"`
		Topics      []string `json:"topics"`
		Constraints []string `json:"constraints"`
		FollowUps   []string `json:"follow_ups"`
		Hints       []string `json:"hints"`
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

func convertStringToDifficulty(diffStr string) (models.ProblemDifficulty, error) {
	switch strings.ToLower(strings.Trim(diffStr, " `")) {
	case "easy":
		return models.EasyDifficulty, nil
	case "medium", "mild":
		return models.MediumDifficulty, nil
	case "hard":
		return models.HardDifficulty, nil
	default:
		return models.EasyDifficulty, fmt.Errorf("'%s' is not a valid difficulty", diffStr)
	}
}

func main() {
	config, err := internal.GetConfig()
	if err != nil {
		panic("failed to load config")
	}

	logger := internal.NewLogger(config.AppEnv)
	defer logger.Sync() //nolint:errcheck // logger.Sync error is not actionable at shutdown

	db, err := internal.GetDBConnFromConfig(&config)
	if err != nil {
		logger.Fatal("failed to connect to database", zap.Error(err))
	}
	defer func() { _ = db.Close() }()

	problemsRepo := repositories.NewProblemPostgresRepository(db, logger)

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

	var mergedProblems problemsResponse
	if err = json.Unmarshal(tagsFileContent, &mergedProblems); err != nil {
		logger.Error("failed to unmarshal problems JSON", zap.Error(err))
		return
	}

	logger.Info("updating", zap.Int("num of problems", len(mergedProblems.Problems)))

	for i := range mergedProblems.Problems {
		problem := &mergedProblems.Problems[i]
		problemDifficulty, err := convertStringToDifficulty(strings.ToLower(problem.Difficulty))
		if err != nil {
			logger.Error("invalid difficulty", zap.String("difficulty", problem.Difficulty), zap.Error(err))
			return
		}

		intProblemID, err := strconv.Atoi(problem.ProblemID)
		if err != nil {
			logger.Error("invalid problem ID", zap.String("problemID", problem.ProblemID), zap.Error(err))
			return
		}

		if err = problemsRepo.SaveProblem(context.Background(), &models.Problem{
			ID:          intProblemID,
			Title:       problem.Title,
			Description: problem.Description,
			Slug:        problem.ProblemSlug,
			Difficulty:  problemDifficulty,
			Topics:      problem.Topics,
		}); err != nil {
			logger.Error("failed to save problem", zap.Int("problemID", intProblemID), zap.Error(err))
			return
		}

		for _, topic := range problem.Topics {
			if err = problemsRepo.SaveProblemTopic(context.Background(), intProblemID, topic); err != nil {
				logger.Error("failed to save problem topic", zap.Int("problemID", intProblemID), zap.String("topic", topic), zap.Error(err))
				return
			}
		}
	}

	logger.Info("data seeding completed successfully")
}
