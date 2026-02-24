package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	fsrs "github.com/open-spaced-repetition/go-fsrs"
)

type ProblemDifficulty string

const (
	EasyDifficulty   ProblemDifficulty = "easy"
	MediumDifficulty ProblemDifficulty = "medium"
	HardDifficulty   ProblemDifficulty = "hard"
)

type ConfidenceLevel string

const (
	AgainConfidence ConfidenceLevel = "again"
	HardConfidence  ConfidenceLevel = "hard"
	GoodConfidence  ConfidenceLevel = "good"
	EasyConfidence  ConfidenceLevel = "easy"
)

type (
	Problem struct {
		ID          int               `json:"id"`
		Topics      []string          `json:"topics"`
		Title       string            `json:"title"`
		Slug        string            `json:"slug"`
		Description string            `json:"description"`
		Difficulty  ProblemDifficulty `json:"difficulty"`
	}

	ProblemTopic struct {
		ID        int    `json:"id"`
		ProblemID int    `json:"problemId"`
		Topic     string `json:"topic"`
	}

	ProblemSubmission struct {
		ID              uuid.UUID       `json:"id"`
		ProblemID       int             `json:"problemId"`
		Date            time.Time       `json:"date"`
		TimeTaken       uint            `json:"timeTaken"`
		ConfidenceLevel ConfidenceLevel `json:"confidenceLevel"`
	}

	ProblemSubmissionWithDetails struct {
		ID              uuid.UUID       `json:"id"`
		SubmittedAt     time.Time       `json:"submittedAt"`
		TimeTaken       uint            `json:"timeTaken"`
		ConfidenceLevel ConfidenceLevel `json:"confidenceLevel"`
		Problem         struct {
			ID          int               `json:"id"`
			Title       string            `json:"title"`
			Description string            `json:"description"`
			Difficulty  ProblemDifficulty `json:"difficulty"`
		} `json:"problem"`
	}

	ProblemSubmissionUserStats struct {
		ID               uuid.UUID     `json:"id"`
		ProblemID        int           `json:"problemID"`
		UserID           uuid.UUID     `json:"userID"`
		NumOfSubmissions uint          `json:"numOfSubmissions"`
		AvgDuration      time.Duration `json:"avgDuration"`
		NextReviewDate   time.Time     `json:"nextReviewDate"`
	}

	ProblemCard struct {
		ProblemID uuid.UUID `json:"problemID"`
		Card      fsrs.Card `json:"card"`
	}
)

func DetermineConfidenceLevelFromString(valStr string) (ConfidenceLevel, error) {
	switch strings.ToLower(strings.TrimSpace(valStr)) {
	case "again", "1":
		return AgainConfidence, nil
	case "hard", "2":
		return HardConfidence, nil
	case "good", "3":
		return GoodConfidence, nil
	case "easy", "4":
		return EasyConfidence, nil
	default:
		return AgainConfidence, fmt.Errorf("%q is not a recognized confidence level", valStr)
	}
}
