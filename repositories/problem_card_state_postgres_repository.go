package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"leetcode-spaced-repetition/models"
)

type ProblemCardStatePostgresRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewProblemCardStatePostgresRepository(db *sql.DB, logger *zap.Logger) *ProblemCardStatePostgresRepository {
	return &ProblemCardStatePostgresRepository{
		db:     db,
		logger: logger,
	}
}

// UpsertProblemCardState implements ProblemCardStateRepository.
func (r *ProblemCardStatePostgresRepository) UpsertProblemCardState(ctx context.Context, state models.ProblemCardState) error {
	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO problem_card_states
			(user_id, problem_id, due, state, stability, difficulty, elapsed_days, scheduled_days, reps, lapses, last_review)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (user_id, problem_id) DO UPDATE SET
			due            = EXCLUDED.due,
			state          = EXCLUDED.state,
			stability      = EXCLUDED.stability,
			difficulty     = EXCLUDED.difficulty,
			elapsed_days   = EXCLUDED.elapsed_days,
			scheduled_days = EXCLUDED.scheduled_days,
			reps           = EXCLUDED.reps,
			lapses         = EXCLUDED.lapses,
			last_review    = EXCLUDED.last_review,
			updated_at     = CURRENT_TIMESTAMP`,
		state.UserID,
		state.ProblemID,
		state.Due,
		state.State,
		state.Stability,
		state.Difficulty,
		state.ElapsedDays,
		state.ScheduledDays,
		state.Reps,
		state.Lapses,
		state.LastReview,
	)
	if err != nil {
		r.logger.Error("failed to upsert problem card state",
			zap.Int("problemID", state.ProblemID),
			zap.String("userID", state.UserID.String()),
			zap.Error(err),
		)
		return err
	}
	return nil
}

// GetProblemCardStateByProblemID implements ProblemCardStateRepository.
func (r *ProblemCardStatePostgresRepository) GetProblemCardStateByProblemID(ctx context.Context, userID uuid.UUID, problemID int) (*models.ProblemCardState, error) {
	var state models.ProblemCardState
	var lastReview sql.NullTime

	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, problem_id, due, state, stability, difficulty, elapsed_days, scheduled_days, reps, lapses, last_review
		FROM problem_card_states
		WHERE user_id = $1 AND problem_id = $2`,
		userID,
		problemID,
	).Scan(
		&state.ID,
		&state.UserID,
		&state.ProblemID,
		&state.Due,
		&state.State,
		&state.Stability,
		&state.Difficulty,
		&state.ElapsedDays,
		&state.ScheduledDays,
		&state.Reps,
		&state.Lapses,
		&lastReview,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		r.logger.Error("failed to get problem card state",
			zap.Int("problemID", problemID),
			zap.String("userID", userID.String()),
			zap.Error(err),
		)
		return nil, err
	}

	if lastReview.Valid {
		state.LastReview = lastReview.Time
	}

	return &state, nil
}

// GetProblemCardStatesDueForReview implements ProblemCardStateRepository.
func (r *ProblemCardStatePostgresRepository) GetProblemCardStatesDueForReview(ctx context.Context, userID uuid.UUID, asOf time.Time, limit int) ([]models.ProblemCardState, error) {
	states := []models.ProblemCardState{}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, problem_id, due, state, stability, difficulty, elapsed_days, scheduled_days, reps, lapses, last_review
		FROM problem_card_states
		WHERE user_id = $1 AND due <= $2
		ORDER BY due ASC
		LIMIT $3`,
		userID,
		asOf,
		limit,
	)
	if err != nil {
		r.logger.Error("failed to query problem card states due for review",
			zap.String("userID", userID.String()),
			zap.Error(err),
		)
		return states, err
	}
	defer rows.Close()

	for rows.Next() {
		var state models.ProblemCardState
		var lastReview sql.NullTime

		if err := rows.Scan(
			&state.ID,
			&state.UserID,
			&state.ProblemID,
			&state.Due,
			&state.State,
			&state.Stability,
			&state.Difficulty,
			&state.ElapsedDays,
			&state.ScheduledDays,
			&state.Reps,
			&state.Lapses,
			&lastReview,
		); err != nil {
			r.logger.Error("failed to scan problem card state row", zap.Error(err))
			return states, err
		}

		if lastReview.Valid {
			state.LastReview = lastReview.Time
		}

		states = append(states, state)
	}

	return states, rows.Err()
}

// GetDueReviewItems implements ProblemCardStateRepository.
func (r *ProblemCardStatePostgresRepository) GetDueReviewItems(ctx context.Context, userID uuid.UUID, asOf time.Time, limit int) ([]models.ProblemReviewItem, error) {
	items := []models.ProblemReviewItem{}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT p.id, p.title, p.slug, p.difficulty, pcs.due, pcs.stability
		FROM problem_card_states pcs
		JOIN problems p ON pcs.problem_id = p.id
		WHERE pcs.user_id = $1 AND pcs.due <= $2
		ORDER BY pcs.due ASC
		LIMIT $3`,
		userID, asOf, limit,
	)
	if err != nil {
		r.logger.Error("failed to query due review items",
			zap.String("userID", userID.String()),
			zap.Error(err),
		)
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.ProblemReviewItem
		if err := rows.Scan(
			&item.ProblemID, &item.Title, &item.Slug,
			&item.Difficulty, &item.Due, &item.Stability,
		); err != nil {
			r.logger.Error("failed to scan due review item row", zap.Error(err))
			return items, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

// GetLowestStabilityItems implements ProblemCardStateRepository.
func (r *ProblemCardStatePostgresRepository) GetLowestStabilityItems(ctx context.Context, userID uuid.UUID, limit int) ([]models.ProblemReviewItem, error) {
	items := []models.ProblemReviewItem{}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT p.id, p.title, p.slug, p.difficulty, pcs.due, pcs.stability
		FROM problem_card_states pcs
		JOIN problems p ON pcs.problem_id = p.id
		WHERE pcs.user_id = $1
		ORDER BY pcs.stability ASC
		LIMIT $2`,
		userID, limit,
	)
	if err != nil {
		r.logger.Error("failed to query lowest stability items",
			zap.String("userID", userID.String()),
			zap.Error(err),
		)
		return items, err
	}
	defer rows.Close()

	for rows.Next() {
		var item models.ProblemReviewItem
		if err := rows.Scan(
			&item.ProblemID, &item.Title, &item.Slug,
			&item.Difficulty, &item.Due, &item.Stability,
		); err != nil {
			r.logger.Error("failed to scan lowest stability item row", zap.Error(err))
			return items, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
