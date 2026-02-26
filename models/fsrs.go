package models

import (
	"time"

	"github.com/google/uuid"
	fsrs "github.com/open-spaced-repetition/go-fsrs/v3"
)

type ProblemCardStateStatus string

const (
	ProblemCardStateNew        ProblemCardStateStatus = "new"
	ProblemCardStateLearning   ProblemCardStateStatus = "learning"
	ProblemCardStateReview     ProblemCardStateStatus = "review"
	ProblemCardStateRelearning ProblemCardStateStatus = "relearning"
)

func ProblemCardStateStatusFromFSRS(s fsrs.State) ProblemCardStateStatus {
	switch s {
	case fsrs.Learning:
		return ProblemCardStateLearning
	case fsrs.Review:
		return ProblemCardStateReview
	case fsrs.Relearning:
		return ProblemCardStateRelearning
	default:
		return ProblemCardStateNew
	}
}

type ProblemCardState struct {
	Due           time.Time
	LastReview    time.Time
	State         ProblemCardStateStatus
	ProblemID     int
	Stability     float64
	Difficulty    float64
	ElapsedDays   uint64
	ScheduledDays uint64
	Reps          uint64
	Lapses        uint64
	ID            uuid.UUID
	UserID        uuid.UUID
}
