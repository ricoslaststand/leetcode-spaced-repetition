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

	"leetcode-spaced-repetition/models"
)

type QuestionPostgresRepository struct {
	db *sql.DB
}

// SaveQuestionSubmission implements QuestionRepository.
func (r QuestionPostgresRepository) SaveQuestionSubmission(c context.Context, questionID int, userID uuid.UUID, date time.Time, timeTaken time.Duration, confidenceLevel models.ConfidenceLevel) error {
	tx, err := r.db.BeginTx(c, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(
		c,
		`INSERT INTO question_submissions (question_id, user_id, submission_date, time_taken, confidence_level) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (question_id, user_id, submission_date) DO NOTHING`,
		questionID,
		userID,
		date,
		fmt.Sprintf("%d seconds", int64(timeTaken.Seconds())),
		confidenceLevel,
	)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23503" {
			return fmt.Errorf("foreign key violation: %v", pqErr.Message)
		}

		return err
	}

	return tx.Commit()
}

// GetAllQuestionsPastReviewDate implements QuestionRepository.
func (r QuestionPostgresRepository) GetAllQuestionsPastReviewDate(c context.Context, limit uint) ([]models.Question, error) {
	var questions []models.Question
	return questions, nil
}

func NewQuestionPostgresRepository(db *sql.DB) *QuestionPostgresRepository {
	return &QuestionPostgresRepository{
		db: db,
	}
}

func (r QuestionPostgresRepository) GetQuestions(ctx context.Context, tags []string, page, limit int) ([]models.Question, error) {
	var questions []models.Question

	rows, err := r.db.QueryContext(
		ctx, `SELECT id, title, slug, difficulty FROM questions WHERE id IN (
		SELECT question_id FROM question_tags WHERE tag IN ($1)
	) ORDER BY id LIMIT $2`, strings.Join(tags, ","), limit)
	if err != nil {
		return questions, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var question models.Question
		err = rows.Scan(&question.ID, &question.Title, &question.Slug, &question.Difficulty)
		if err != nil {
			return questions, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

func (r QuestionPostgresRepository) GetQuestionByID(ctx context.Context, questionID int) (*models.Question, error) {
	var id int
	var title string
	var slug string
	var difficulty models.QuestionDifficulty

	row := r.db.QueryRowContext(ctx, "SELECT id, title, slug, difficulty FROM questions WHERE id = $1", questionID)
	switch err := row.Scan(&id, &title, &slug, &difficulty); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &models.Question{
			ID:         id,
			Title:      title,
			Slug:       slug,
			Difficulty: difficulty,
		}, nil
	default:
		return nil, err
	}
}

func (r QuestionPostgresRepository) GetQuestionStatsByID(ctx context.Context, id int) (*models.QuestionSubmissionUserStats, error) {
	return nil, nil
}

func (r QuestionPostgresRepository) GetQuestionSubmissions(ctx context.Context, questionID []int) ([]models.QuestionSubmissionWithDetails, error) {
	var submissions []models.QuestionSubmissionWithDetails
	if len(questionID) == 0 {
		rows, err := r.db.QueryContext(
			ctx,
			`SELECT
				question_submissions.id,
				question_submissions.question_id,
				submission_date,
				EXTRACT(EPOCH  FROM time_taken),
				confidence_level,
				questions.title,
				questions.description,
				questions.difficulty
			FROM question_submissions
			JOIN questions ON question_submissions.question_id = questions.id
			ORDER BY submission_date DESC`,
		)
		if err != nil {
			return []models.QuestionSubmissionWithDetails{}, err
		}
		defer func() { _ = rows.Close() }()

		for rows.Next() {
			var sub models.QuestionSubmissionWithDetails

			if err := rows.Scan(
				&sub.ID,
				&sub.Question.ID,
				&sub.SubmittedAt,
				&sub.TimeTaken,
				&sub.ConfidenceLevel,
				&sub.Question.Title,
				&sub.Question.Description,
				&sub.Question.Difficulty,
			); err != nil {
				return []models.QuestionSubmissionWithDetails{}, err
			}

			submissions = append(submissions, sub)
		}

		return submissions, nil
	} else {
		rows, err := r.db.QueryContext(
			ctx,
			`SELECT
				question_submissions.id,
				question_submissions.question_id,
				submission_date,
				EXTRACT(EPOCH  FROM time_taken),
				confidence_level,
				questions.title,
				questions.description
			FROM question_submissions
			JOIN questions ON question_submissions.question_id = questions.id WHERE question_id = ANY($1)
			ORDER BY submission_date DESC`,
			pq.Array(questionID),
		)
		if err != nil {
			return []models.QuestionSubmissionWithDetails{}, err
		}
		defer func() { _ = rows.Close() }()

		for rows.Next() {
			var sub models.QuestionSubmissionWithDetails

			if err := rows.Scan(
				&sub.ID,
				&sub.Question.ID,
				&sub.SubmittedAt,
				&sub.TimeTaken,
				&sub.ConfidenceLevel,
				&sub.Question.Title,
				&sub.Question.Description,
			); err != nil {
				return []models.QuestionSubmissionWithDetails{}, err
			}

			submissions = append(submissions, sub)
		}

		return submissions, nil
	}
}

func (r QuestionPostgresRepository) GetSubmissionsByQuestionID(c context.Context, questionID int) ([]models.QuestionSubmission, error) {
	var submissions []models.QuestionSubmission

	rows, err := r.db.QueryContext(
		c,
		`SELECT
			id,
			question_id,
			submission_date,
			EXTRACT(EPOCH  FROM time_taken),
			confidence_level
		FROM question_submissions
		WHERE question_id = $1
		ORDER BY submission_date DESC`,
		questionID,
	)
	if err != nil {
		return []models.QuestionSubmission{}, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var sub models.QuestionSubmission

		if err := rows.Scan(
			&sub.ID,
			&sub.QuestionID,
			&sub.Date,
			&sub.TimeTaken,
			&sub.ConfidenceLevel,
		); err != nil {
			return []models.QuestionSubmission{}, err
		}

		submissions = append(submissions, sub)
	}

	return submissions, nil
}

func (r QuestionPostgresRepository) SaveQuestion(c context.Context, q *models.Question) error {
	_, err := r.db.Exec(
		"INSERT INTO questions (id, title, slug, description, difficulty) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING",
		q.ID, q.Title, q.Slug, q.Description, q.Difficulty,
	)

	return err
}

func (r QuestionPostgresRepository) SaveQuestionTag(c context.Context, questionId int, tag string) error {
	_, err := r.db.Exec(
		"INSERT INTO question_tags (question_id, tag) VALUES ($1, $2) ON CONFLICT (question_id, tag) DO NOTHING",
		questionId, tag,
	)

	return err
}

func (r QuestionPostgresRepository) GetAllQuestionTags(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT DISTINCT(tag) FROM question_tags ORDER BY tag")
	if err != nil {
		return []string{}, err
	}
	defer func() { _ = rows.Close() }()

	var tags []string

	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return []string{}, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r QuestionPostgresRepository) GetTagsForQuestion(ctx context.Context, id int) ([]string, error) {
	rows, err := r.db.QueryContext(ctx, "SELECT tag FROM question_tags WHERE question_id = $1", id)
	if err != nil {
		return []string{}, err
	}
	defer func() { _ = rows.Close() }()

	var tags []string

	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return []string{}, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
