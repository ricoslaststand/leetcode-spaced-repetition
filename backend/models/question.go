package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type QuestionDifficulty int

const (
	EasyDifficulty QuestionDifficulty = iota + 1
	MediumDifficulty
	HardDifficulty
)

type ConfidenceLevel int

const (
	VeryLowConfidence ConfidenceLevel = iota + 1
	LowConfidence
	MediumConfidence
	HighConfidence
	VeryHighConfidence
)

type Question struct {
	ID          int                `json:"id"`
	Tags        []string           `json:"tags"`
	Title       string             `json:"title"`
	Slug        string             `json:"slug"`
	Description string             `json:"description"`
	Difficulty  QuestionDifficulty `json:"difficulty"`
}

type QuestionTag struct {
	ID         int    `json:"id"`
	QuestionID int    `json:"questionId"`
	Tag        string `json:"tag"`
}

type QuestionSubmission struct {
	ID              uuid.UUID       `json:"id"`
	QuestionID      int             `json:"questionId"`
	Date            time.Time       `json:"date"`
	ConfidenceLevel ConfidenceLevel `json:"confidenceLevel"`
}

type QuestionSubmissionUserStats struct {
	ID               uuid.UUID     `json:"id"`
	QuestionID       int           `json:"questionID"`
	UserID           uuid.UUID     `json:"userID"`
	NumOfSubmissions uint          `json:"numOfSubmissions"`
	AvgDuration      time.Duration `json:"avgDuration"`
}

func DetermineDifficulty(val int) (QuestionDifficulty, error) {
	if val < int(EasyDifficulty) && val > int(HardDifficulty) {
		return EasyDifficulty, fmt.Errorf("%d is not recognized as a valid difficulty level", val)
	}

	return QuestionDifficulty(val), nil
}
