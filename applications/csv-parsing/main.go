package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"leetcode-spaced-repetition/internal"
	"leetcode-spaced-repetition/models"
	"leetcode-spaced-repetition/repositories"
	"leetcode-spaced-repetition/services"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

const (
	ProblemNumberIdx int = iota
	DateIdx
	TimeTakenIdx
	ConfidenceLevelIdx
)

const dateFormat string = "2006-01-02"

type recordSubmission struct {
	problemNumber   int
	timeTaken       *time.Duration
	submissionDate  time.Time
	confidenceLevel models.ConfidenceLevel
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

	logger.Info("constructing domain layers")

	problemsRepo := repositories.NewProblemPostgresRepository(db, logger)
	problemsService := services.NewProblemsService(problemsRepo, logger)

	file, err := os.Open("leetcode_submissions.csv")
	if err != nil {
		logger.Fatal("failed to open CSV file", zap.Error(err))
	}
	defer file.Close()

	logger.Info("reading leetcode submission file")

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		logger.Fatal("failed to read CSV file", zap.Error(err))
	}

	logger.Info("processing records", zap.Int("count", len(records)))

	c := context.Background()

	for i := 1; i < len(records); i++ {
		logger.Debug("processing record", zap.Int("index", i), zap.Strings("record", records[i]))
		submission, err := validateProblemSubmission(records[i])
		if err != nil {
			logger.Error("invalid submission record", zap.Int("index", i), zap.Error(err))
			continue
		}

		userID := uuid.MustParse("a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11")

		if err = problemsService.SaveProblemSubmission(c, submission.problemNumber, userID, submission.submissionDate, submission.timeTaken, submission.confidenceLevel); err != nil {
			logger.Error("failed to save problem submission", zap.Int("problemNumber", submission.problemNumber), zap.Error(err))
		} else {
			logger.Info("saved problem submission", zap.Int("problemNumber", submission.problemNumber))
		}
	}

	logger.Info("CSV parsing completed successfully")
}

func validateProblemSubmission(r []string) (recordSubmission, error) {
	problemNum := r[ProblemNumberIdx]
	date := r[DateIdx]
	timeTaken := r[TimeTakenIdx]
	confidenceLevelStr := r[ConfidenceLevelIdx]

	problemNumber, err := strconv.ParseInt(problemNum, 10, 0)
	if err != nil {
		return recordSubmission{}, fmt.Errorf("'%s' is not a valid problem number", problemNum)
	}

	var timeTakenDuration *time.Duration
	formattedTimeTakenStr := strings.ReplaceAll(strings.TrimSpace(timeTaken), " ", "")
	if formattedTimeTakenStr != "" {
		d, err := time.ParseDuration(formattedTimeTakenStr)
		if err != nil {
			return recordSubmission{}, fmt.Errorf("'%s' is not a valid time duration", formattedTimeTakenStr)
		}
		timeTakenDuration = &d
	}

	dateTime, err := time.Parse(dateFormat, date)
	if err != nil {
		return recordSubmission{}, fmt.Errorf("'%s' is not a valid date format", date)
	}

	confidenceLevel, err := models.DetermineConfidenceLevelFromString(confidenceLevelStr)
	if err != nil {
		return recordSubmission{}, fmt.Errorf("'%s' is not a valid confidence level", confidenceLevelStr)
	}

	return recordSubmission{
		problemNumber:   int(problemNumber),
		timeTaken:       timeTakenDuration,
		submissionDate:  dateTime,
		confidenceLevel: confidenceLevel,
	}, nil

}
