package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	fsrs "github.com/open-spaced-repetition/go-fsrs"
)

type QuestionDifficulty string

const (
	EasyDifficulty QuestionDifficulty = "easy"
	MildDifficulty QuestionDifficulty = "mild"
	HardDifficulty QuestionDifficulty = "hard"
)

type ConfidenceLevel string

const (
	AgainConfidence ConfidenceLevel = "again"
	HardConfidence  ConfidenceLevel = "hard"
	GoodConfidence  ConfidenceLevel = "good"
	EasyConfidence  ConfidenceLevel = "easy"
)

type (
	Question struct {
		ID          int                `json:"id"`
		Tags        []string           `json:"tags"`
		Title       string             `json:"title"`
		Slug        string             `json:"slug"`
		Description string             `json:"description"`
		Difficulty  QuestionDifficulty `json:"difficulty"`
	}

	QuestionTag struct {
		ID         int    `json:"id"`
		QuestionID int    `json:"questionId"`
		Tag        string `json:"tag"`
	}

	QuestionSubmission struct {
		ID              uuid.UUID       `json:"id"`
		QuestionID      int             `json:"questionId"`
		Date            time.Time       `json:"date"`
		TimeTaken       uint            `json:"timeTaken"`
		ConfidenceLevel ConfidenceLevel `json:"confidenceLevel"`
	}

	QuestionSubmissionWithDetails struct {
		ID              uuid.UUID       `json:"id"`
		SubmittedAt     time.Time       `json:"submittedAt"`
		TimeTaken       uint            `json:"timeTaken"`
		ConfidenceLevel ConfidenceLevel `json:"confidenceLevel"`
		Question        struct {
			ID          int                `json:"id"`
			Title       string             `json:"title"`
			Description string             `json:"description"`
			Difficulty  QuestionDifficulty `json:"difficulty"`
		} `json:"question"`
	}

	QuestionSubmissionUserStats struct {
		ID               uuid.UUID     `json:"id"`
		QuestionID       int           `json:"questionID"`
		UserID           uuid.UUID     `json:"userID"`
		NumOfSubmissions uint          `json:"numOfSubmissions"`
		AvgDuration      time.Duration `json:"avgDuration"`
		NextReviewDate   time.Time     `json:"nextReviewDate"`
	}

	QuestionCard struct {
		QuestionID uuid.UUID `json:"questionID"`
		Card       fsrs.Card `json:"card"`
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
