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
	"github.com/mgechev/revive/config"
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
	questionNumber  int
	timeTaken       time.Duration
	submissionDate  time.Time
	confidenceLevel models.ConfidenceLevel
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

	logger.Info("constructing domain layers")

	questionsRepo := repositories.NewQuestionPostgresRepository(db, logger)
	questionsService := services.NewQuestionsService(questionsRepo, logger)

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
		submission, err := validateQuestionSubmission(records[i])
		if err != nil {
			logger.Error("invalid submission record", zap.Int("index", i), zap.Error(err))
			continue
		}

		userID, _ := uuid.NewUUID()

		if err = questionsService.SaveQuestionSubmission(c, submission.questionNumber, userID, submission.submissionDate, submission.timeTaken, submission.confidenceLevel); err != nil {
			logger.Error("failed to save question submission", zap.Int("questionNumber", submission.questionNumber), zap.Error(err))
		} else {
			logger.Info("saved question submission", zap.Int("questionNumber", submission.questionNumber))
		}
	}

	logger.Info("CSV parsing completed successfully")
}

func validateQuestionSubmission(r []string) (recordSubmission, error) {
	problemNum := r[ProblemNumberIdx]
	date := r[DateIdx]
	timeTaken := r[TimeTakenIdx]
	confidenceLevelStr := r[ConfidenceLevelIdx]

	questionNum, err := strconv.ParseInt(problemNum, 10, 0)
	if err != nil {
		return recordSubmission{}, fmt.Errorf("'%s' is not a valid question number", problemNum)
	}

	formattedTimeTakenDuration := strings.ReplaceAll(timeTaken, " ", "")
	timeTakenDuration, err := time.ParseDuration(formattedTimeTakenDuration)
	if err != nil {
		return recordSubmission{}, fmt.Errorf("'%s' is not a valid time duration", formattedTimeTakenDuration)
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
		questionNumber:  int(questionNum),
		timeTaken:       timeTakenDuration,
		submissionDate:  dateTime,
		confidenceLevel: confidenceLevel,
	}, nil
}
