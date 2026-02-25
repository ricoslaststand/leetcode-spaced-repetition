package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.uber.org/zap"

	"leetcode-spaced-repetition/models"
)

type ProblemPostgresRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewProblemPostgresRepository(db *sql.DB, logger *zap.Logger) *ProblemPostgresRepository {
	return &ProblemPostgresRepository{
		db:     db,
		logger: logger,
	}
}

// SaveProblemSubmission implements ProblemRepository.
// Returns (true, nil) if a row was inserted, (false, nil) if skipped due to a conflict, or (false, err) on error.
func (r ProblemPostgresRepository) SaveProblemSubmission(c context.Context, problemID int, userID uuid.UUID, date time.Time, timeTaken *time.Duration, confidenceLevel models.ConfidenceLevel, language *models.SubmissionLanguage) (bool, error) {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return false, err
	}
	defer func() { _ = tx.Rollback() }()

	var timeTakenVal any
	if timeTaken != nil {
		timeTakenVal = fmt.Sprintf("%d seconds", int64(timeTaken.Seconds()))
	}

	result, err := tx.ExecContext(
		c,
		`INSERT INTO problem_submissions (problem_id, user_id, submission_date, time_taken, confidence_level, language) VALUES ($1, $2, $3, $4, $5, $6) ON CONFLICT (problem_id, user_id, submission_date) DO NOTHING`,
		problemID,
		userID,
		date,
		timeTakenVal,
		confidenceLevel,
		language,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			r.logger.Error("foreign key violation saving problem submission", zap.Int("problemID", problemID), zap.Error(err))
			return false, fmt.Errorf("foreign key violation: %v", pqErr.Message)
		}
		r.logger.Error("failed to save problem submission", zap.Int("problemID", problemID), zap.Error(err))
		return false, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("failed to get rows affected", zap.Int("problemID", problemID), zap.Error(err))
		return false, err
	}

	if err := tx.Commit(); err != nil {
		return false, err
	}

	return rowsAffected == 1, nil
}

// GetAllProblemsPastReviewDate implements ProblemRepository.
func (r ProblemPostgresRepository) GetAllProblemsPastReviewDate(c context.Context, limit uint) ([]models.Problem, error) {
	var problems []models.Problem
	return problems, nil
}

func (r ProblemPostgresRepository) GetProblems(ctx context.Context, topics []string, page, limit int) ([]models.Problem, error) {
	var problems []models.Problem

	rows, err := r.db.QueryContext(
		ctx, `SELECT id, title, slug, difficulty FROM problems WHERE id IN (
		SELECT problem_id FROM problem_topics WHERE topic IN ($1)
	) ORDER BY id LIMIT $2`, strings.Join(topics, ","), limit)
	if err != nil {
		r.logger.Error("failed to query problems", zap.Error(err))
		return problems, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var problem models.Problem
		err = rows.Scan(&problem.ID, &problem.Title, &problem.Slug, &problem.Difficulty)
		if err != nil {
			r.logger.Error("failed to scan problem row", zap.Error(err))
			return problems, err
		}
		problems = append(problems, problem)
	}

	return problems, nil
}

func (r ProblemPostgresRepository) GetProblemByID(ctx context.Context, problemID int) (*models.Problem, error) {
	var id int
	var title string
	var slug string
	var difficulty models.ProblemDifficulty

	row := r.db.QueryRowContext(ctx, "SELECT id, title, slug, difficulty FROM problems WHERE id = $1", problemID)
	switch err := row.Scan(&id, &title, &slug, &difficulty); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &models.Problem{
			ID:         id,
			Title:      title,
			Slug:       slug,
			Difficulty: difficulty,
		}, nil
	default:
		r.logger.Error("failed to get problem by ID", zap.Int("problemID", problemID), zap.Error(err))
		return nil, err
	}
}

func (r ProblemPostgresRepository) GetProblemStatsByID(ctx context.Context, id int) (*models.ProblemSubmissionUserStats, error) {
	return nil, nil
}

func (r ProblemPostgresRepository) GetProblemSubmissions(ctx context.Context, problemID []int) ([]models.ProblemSubmissionWithDetails, error) {
	var submissions []models.ProblemSubmissionWithDetails
	if len(problemID) == 0 {
		rows, err := r.db.QueryContext(
			ctx,
			`SELECT
				problem_submissions.id,
				problem_submissions.problem_id,
				submission_date,
				EXTRACT(EPOCH  FROM time_taken),
				confidence_level,
				problems.title,
				problems.description,
				problems.difficulty
			FROM problem_submissions
			JOIN problems ON problem_submissions.problem_id = problems.id
			ORDER BY submission_date DESC`,
		)
		if err != nil {
			r.logger.Error("failed to query all problem submissions", zap.Error(err))
			return []models.ProblemSubmissionWithDetails{}, err
		}
		defer func() { _ = rows.Close() }()

		for rows.Next() {
			var sub models.ProblemSubmissionWithDetails
			var timeTakenSeconds sql.NullFloat64

			if err := rows.Scan(
				&sub.ID,
				&sub.Problem.ID,
				&sub.SubmittedAt,
				&timeTakenSeconds,
				&sub.ConfidenceLevel,
				&sub.Problem.Title,
				&sub.Problem.Description,
				&sub.Problem.Difficulty,
			); err != nil {
				r.logger.Error("failed to scan problem submission row", zap.Error(err))
				return []models.ProblemSubmissionWithDetails{}, err
			}

			if timeTakenSeconds.Valid {
				v := uint(timeTakenSeconds.Float64)
				sub.TimeTaken = &v
			}

			submissions = append(submissions, sub)
		}

		return submissions, nil
	} else {
		rows, err := r.db.QueryContext(
			ctx,
			`SELECT
				problem_submissions.id,
				problem_submissions.problem_id,
				submission_date,
				EXTRACT(EPOCH  FROM time_taken),
				confidence_level,
				problems.title,
				problems.description
			FROM problem_submissions
			JOIN problems ON problem_submissions.problem_id = problems.id WHERE problem_id = ANY($1)
			ORDER BY submission_date DESC`,
			pq.Array(problemID),
		)
		if err != nil {
			r.logger.Error("failed to query problem submissions by IDs", zap.Error(err))
			return []models.ProblemSubmissionWithDetails{}, err
		}
		defer func() { _ = rows.Close() }()

		for rows.Next() {
			var sub models.ProblemSubmissionWithDetails
			var timeTakenSeconds sql.NullFloat64

			if err := rows.Scan(
				&sub.ID,
				&sub.Problem.ID,
				&sub.SubmittedAt,
				&timeTakenSeconds,
				&sub.ConfidenceLevel,
				&sub.Problem.Title,
				&sub.Problem.Description,
			); err != nil {
				r.logger.Error("failed to scan problem submission row", zap.Error(err))
				return []models.ProblemSubmissionWithDetails{}, err
			}

			if timeTakenSeconds.Valid {
				v := uint(timeTakenSeconds.Float64)
				sub.TimeTaken = &v
			}

			submissions = append(submissions, sub)
		}

		return submissions, nil
	}
}

func (r ProblemPostgresRepository) GetSubmissionsByProblemID(c context.Context, problemID int) ([]models.ProblemSubmission, error) {
	var submissions []models.ProblemSubmission

	rows, err := r.db.QueryContext(
		c,
		`SELECT
			id,
			problem_id,
			submission_date,
			EXTRACT(EPOCH  FROM time_taken),
			confidence_level
		FROM problem_submissions
		WHERE problem_id = $1
		ORDER BY submission_date DESC`,
		problemID,
	)
	if err != nil {
		r.logger.Error("failed to query submissions by problem ID", zap.Int("problemID", problemID), zap.Error(err))
		return []models.ProblemSubmission{}, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var sub models.ProblemSubmission
		var timeTakenSeconds sql.NullFloat64

		if err := rows.Scan(
			&sub.ID,
			&sub.ProblemID,
			&sub.Date,
			&timeTakenSeconds,
			&sub.ConfidenceLevel,
		); err != nil {
			r.logger.Error("failed to scan submission row", zap.Error(err))
			return []models.ProblemSubmission{}, err
		}

		if timeTakenSeconds.Valid {
			v := uint(timeTakenSeconds.Float64)
			sub.TimeTaken = &v
		}

		submissions = append(submissions, sub)
	}

	return submissions, nil
}

func (r ProblemPostgresRepository) SaveProblem(c context.Context, p *models.Problem) error {
	_, err := r.db.Exec(
		"INSERT INTO problems (id, title, slug, description, difficulty) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING",
		p.ID, p.Title, p.Slug, p.Description, p.Difficulty,
	)
	if err != nil {
		r.logger.Error("failed to save problem", zap.Int("problemID", p.ID), zap.Error(err))
	}

	return err
}

func (r ProblemPostgresRepository) SaveProblemTopic(c context.Context, problemId int, topic string) error {
	_, err := r.db.Exec(
		"INSERT INTO problem_topics (problem_id, topic) VALUES ($1, $2) ON CONFLICT (problem_id, topic) DO NOTHING",
		problemId, topic,
	)
	if err != nil {
		r.logger.Error("failed to save problem topic", zap.Int("problemID", problemId), zap.String("topic", topic), zap.Error(err))
	}

	return err
}

func (r ProblemPostgresRepository) GetAllProblemTopics(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT DISTINCT(topic) FROM problem_topics ORDER BY topic")
	if err != nil {
		r.logger.Error("failed to query all problem topics", zap.Error(err))
		return []string{}, err
	}
	defer func() { _ = rows.Close() }()

	var topics []string

	for rows.Next() {
		var topic string
		err = rows.Scan(&topic)
		if err != nil {
			r.logger.Error("failed to scan topic row", zap.Error(err))
			return []string{}, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
}

func (r ProblemPostgresRepository) GetTopicsForProblem(ctx context.Context, id int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT topic FROM problem_topics WHERE problem_id = $1", id)
	if err != nil {
		r.logger.Error("failed to query topics for problem", zap.Int("problemID", id), zap.Error(err))
		return []string{}, err
	}
	defer func() { _ = rows.Close() }()

	var topics []string

	for rows.Next() {
		var topic string
		err = rows.Scan(&topic)
		if err != nil {
			r.logger.Error("failed to scan topic row", zap.Error(err))
			return []string{}, err
		}
		topics = append(topics, topic)
	}

	return topics, nil
}
