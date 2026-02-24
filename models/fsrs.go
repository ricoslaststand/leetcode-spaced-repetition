package models

import (
	"time"

	"github.com/google/uuid"
)

type CardState struct {
	ID            uuid.UUID
	ProblemID     int
	UserID        uuid.UUID
	Stability     float64
	Difficulty    float64
	ElapsedDays   uint64
	ScheduledDays uint64
	Reps          uint64
	Lapses        uint64
	LastReview    time.Time
}
