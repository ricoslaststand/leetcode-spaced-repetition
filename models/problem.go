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

type SubmissionLanguage string

const (
	LangCPP        SubmissionLanguage = "cpp"
	LangJava       SubmissionLanguage = "java"
	LangPython     SubmissionLanguage = "python"
	LangPython3    SubmissionLanguage = "python3"
	LangC          SubmissionLanguage = "c"
	LangCSharp     SubmissionLanguage = "csharp"
	LangJavaScript SubmissionLanguage = "javascript"
	LangTypeScript SubmissionLanguage = "typescript"
	LangPHP        SubmissionLanguage = "php"
	LangSwift      SubmissionLanguage = "swift"
	LangKotlin     SubmissionLanguage = "kotlin"
	LangDart       SubmissionLanguage = "dart"
	LangGolang     SubmissionLanguage = "golang"
	LangRuby       SubmissionLanguage = "ruby"
	LangScala      SubmissionLanguage = "scala"
	LangRust       SubmissionLanguage = "rust"
	LangRacket     SubmissionLanguage = "racket"
	LangErlang     SubmissionLanguage = "erlang"
	LangElixir     SubmissionLanguage = "elixir"
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
		ID              uuid.UUID          `json:"id"`
		ProblemID       int                `json:"problemId"`
		Date            time.Time          `json:"date"`
		TimeTaken       *uint              `json:"timeTaken"`
		ConfidenceLevel ConfidenceLevel    `json:"confidenceLevel"`
		Language        *SubmissionLanguage `json:"language"`
	}

	ProblemSubmissionWithDetails struct {
		ID              uuid.UUID       `json:"id"`
		SubmittedAt     time.Time       `json:"submittedAt"`
		TimeTaken       *uint           `json:"timeTaken"`
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

type ProblemDetail struct {
	Problem
	Stability *float64 `json:"stability"`
}

type ProblemReviewItem struct {
	ProblemID  int               `json:"problemId"`
	Title      string            `json:"title"`
	Slug       string            `json:"slug"`
	Difficulty ProblemDifficulty `json:"difficulty"`
	Due        time.Time         `json:"due"`
	Stability  float64           `json:"stability"`
}

type DashboardData struct {
	Due          []ProblemReviewItem `json:"due"`
	LowStability []ProblemReviewItem `json:"lowStability"`
}

type ImportSubmissionRowError struct {
	Row    int    `json:"row"`
	Reason string `json:"reason"`
}

type ImportSubmissionsResult struct {
	Imported int                        `json:"imported"`
	Skipped  int                        `json:"skipped"`
	Errors   []ImportSubmissionRowError `json:"errors"`
}

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
