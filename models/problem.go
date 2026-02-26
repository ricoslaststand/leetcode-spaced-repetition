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
		Title       string            `json:"title"`
		Slug        string            `json:"slug"`
		Description string            `json:"description"`
		Difficulty  ProblemDifficulty `json:"difficulty"`
		Topics      []string          `json:"topics"`
		ID          int               `json:"id"`
	}

	ProblemTopic struct {
		Topic     string `json:"topic"`
		ID        int    `json:"id"`
		ProblemID int    `json:"problemId"`
	}

	ProblemSubmission struct {
		Date            time.Time           `json:"date"`
		TimeTaken       *uint               `json:"timeTaken"`
		Language        *SubmissionLanguage `json:"language"`
		ConfidenceLevel ConfidenceLevel     `json:"confidenceLevel"`
		ProblemID       int                 `json:"problemId"`
		ID              uuid.UUID           `json:"id"`
	}

	ProblemSubmissionWithDetails struct {
		SubmittedAt     time.Time       `json:"submittedAt"`
		TimeTaken       *uint           `json:"timeTaken"`
		ConfidenceLevel ConfidenceLevel `json:"confidenceLevel"`
		Problem         struct {
			Title       string            `json:"title"`
			Description string            `json:"description"`
			Difficulty  ProblemDifficulty `json:"difficulty"`
			ID          int               `json:"id"`
		} `json:"problem"`
		ID uuid.UUID `json:"id"`
	}

	ProblemSubmissionUserStats struct {
		NextReviewDate   time.Time     `json:"nextReviewDate"`
		ProblemID        int           `json:"problemID"`
		NumOfSubmissions uint          `json:"numOfSubmissions"`
		AvgDuration      time.Duration `json:"avgDuration"`
		ID               uuid.UUID     `json:"id"`
		UserID           uuid.UUID     `json:"userID"`
	}

	ProblemCard struct {
		Card      fsrs.Card `json:"card"`
		ProblemID uuid.UUID `json:"problemID"`
	}
)

type ProblemDetail struct {
	Stability *float64 `json:"stability"`
	Problem
}

type ProblemReviewItem struct {
	Due        time.Time         `json:"due"`
	Title      string            `json:"title"`
	Slug       string            `json:"slug"`
	Difficulty ProblemDifficulty `json:"difficulty"`
	ProblemID  int               `json:"problemId"`
	Stability  float64           `json:"stability"`
}

type DashboardData struct {
	Due          []ProblemReviewItem `json:"due"`
	LowStability []ProblemReviewItem `json:"lowStability"`
	OverdueCount int                 `json:"overdueCount"`
}

type ImportSubmissionRowError struct {
	Reason string `json:"reason"`
	Row    int    `json:"row"`
}

type ImportSubmissionsResult struct {
	Errors   []ImportSubmissionRowError `json:"errors"`
	Imported int                        `json:"imported"`
	Skipped  int                        `json:"skipped"`
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
