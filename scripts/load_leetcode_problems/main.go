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

type MergedProblems struct {
	Problems []struct {
		ID     string   `json:"problem_id"`
		Topics []string `json:"topics"`
	} `json:"questions"`
}

func main() {
	cfg, err := internal.GetConfig()
	if err != nil {
		panic("failed to load config")
	}

	logger := internal.NewLogger(cfg.AppEnv)
	defer logger.Sync() //nolint:errcheck

	db, err := internal.GetDBConnFromConfig(cfg)
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

	numOfProblems := len(responseData.StatStatusPairs)
	logger.Info("downloaded problems", zap.Int("count", numOfProblems))

	problemRepo := repositories.NewProblemPostgresRepository(db, logger)

	var problems []models.Problem
	for i := 0; i < numOfProblems; i++ {
		currProblem := responseData.StatStatusPairs[i]
		difficultyMap := map[int]models.ProblemDifficulty{
			1: models.EasyDifficulty,
			2: models.MediumDifficulty,
			3: models.HardDifficulty,
		}
		problemDifficulty, ok := difficultyMap[currProblem.Difficulty.Level]
		if !ok {
			logger.Error("unrecognized difficulty level", zap.Int("level", currProblem.Difficulty.Level))
			return
		}

		problems = append(problems, models.Problem{
			ID:          currProblem.Stat.QuestionID,
			Title:       currProblem.Stat.QuestionTitle,
			Description: "",
			Slug:        currProblem.Stat.QuestionTitleSlug,
			Difficulty:  problemDifficulty,
		})
	}

	logger.Info("saving problems to the database", zap.Int("count", len(problems)))

	c := context.Background()

	for i := 0; i < len(problems); i++ {
		if err = problemRepo.SaveProblem(c, &problems[i]); err != nil {
			logger.Fatal("failed to save problem", zap.Int("problemID", problems[i].ID), zap.Error(err))
		}
	}

	// Adding problem topics
	tagsFileContent, err := os.ReadFile("merged_problems.json")
	if err != nil {
		logger.Fatal("failed to read merged_problems.json", zap.Error(err))
	}

	var mergedProblems MergedProblems
	if err = json.Unmarshal(tagsFileContent, &mergedProblems); err != nil {
		logger.Fatal("failed to unmarshal topics JSON", zap.Error(err))
	}

	for i := 0; i < len(mergedProblems.Problems); i++ {
		problem := mergedProblems.Problems[i]
		problemId, err := strconv.Atoi(problem.ID)
		if err != nil {
			logger.Error("invalid problem ID", zap.String("id", problem.ID), zap.Error(err))
			continue
		}

		for j := 0; j < len(problem.Topics); j++ {
			logger.Debug("processing topic", zap.Int("problemID", problemId), zap.String("topic", problem.Topics[j]))
			if err = problemRepo.SaveProblemTopic(c, problemId, problem.Topics[j]); err != nil {
				logger.Fatal("failed to save problem topic", zap.Int("problemID", problemId), zap.String("topic", problem.Topics[j]), zap.Error(err))
			}
		}
	}

	logger.Info("load_leetcode_problems completed successfully")
	fmt.Println()
}
