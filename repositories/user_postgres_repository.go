package repositories

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"leetcode-spaced-repetition/models"
)

type UserPostgresRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewUserPostgresRepository(db *sql.DB, logger *zap.Logger) *UserPostgresRepository {
	return &UserPostgresRepository{db: db, logger: logger}
}

func (r UserPostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User

	if err := r.db.QueryRowContext(
		ctx,
		"SELECT FROM id, email FROM users WHERE email = $1",
		email,
	).Scan(&user.ID, &user.Email); err != nil {
		r.logger.Error("failed to get user by email", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (r UserPostgresRepository) GetPasswordHashByEmail(ctx context.Context, email string) (*string, error) {
	var passwordHash string

	if err := r.db.QueryRowContext(
		ctx,
		"SELECT passwordHash FROM users WHERE email = $1",
		email,
	).Scan(&passwordHash); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("failed to get password hash by email", zap.String("email", email), zap.Error(err))
		return nil, err
	}

	return &passwordHash, nil
}
