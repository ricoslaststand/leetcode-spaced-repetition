package repositories

import (
	"context"
	"time"

	"github.com/google/uuid"

	"leetcode-spaced-repetition/models"
)

type ProblemCardStateRepository interface {
	// UpsertProblemCardState inserts or updates the FSRS state for a (user, problem) pair.
	UpsertProblemCardState(ctx context.Context, state models.ProblemCardState) error

	// GetProblemCardStateByProblemID returns the current FSRS state for a (user, problem) pair,
	// or nil if no state has been recorded yet.
	GetProblemCardStateByProblemID(ctx context.Context, userID uuid.UUID, problemID int) (*models.ProblemCardState, error)

	// GetProblemCardStatesDueForReview returns cards whose due date is on or before asOf.
	GetProblemCardStatesDueForReview(ctx context.Context, userID uuid.UUID, asOf time.Time, limit int) ([]models.ProblemCardState, error)
}
