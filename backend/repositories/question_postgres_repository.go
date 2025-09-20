package repositories

import (
	"database/sql"
	"leetcode-spaced-repetition/models"
)

type QuestionPostgresRepository struct {
	db *sql.DB
}

func NewQuestionPostgresRepository(db *sql.DB) *QuestionPostgresRepository {
	return &QuestionPostgresRepository{
		db: db,
	}
}

func (r QuestionPostgresRepository) GetQuestions() ([]models.Question, error) {
	return []models.Question{}, nil
}

func (r QuestionPostgresRepository) GetQuestionByID(ID int) (*models.Question, error) {
	var id int
	var title string
	var slug string
	var difficulty int

	row := r.db.QueryRow("SELECT id, title, slug, difficulty FROM questions WHERE id = $1", ID)
	switch err := row.Scan(&id, &title, &slug, &difficulty); err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return &models.Question{
			ID:         id,
			Title:      title,
			Slug:       slug,
			Difficulty: models.QuestionDifficulty(difficulty),
		}, nil
	default:
		return nil, err
	}
}

func (r QuestionPostgresRepository) GetQuestionSubmissions() ([]models.QuestionSubmission, error) {
	var submissions []models.QuestionSubmission

	return submissions, nil
}

func (r QuestionPostgresRepository) SaveQuestion(q models.Question) error {
	_, err := r.db.Exec(
		"INSERT INTO questions (id, title, slug, description, difficulty) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (id) DO NOTHING",
		q.ID, q.Title, q.Slug, q.Description, q.Difficulty,
	)

	return err
}

func (r QuestionPostgresRepository) SaveQuestionTag(questionId int, tag string) error {
	_, err := r.db.Exec(
		"INSERT INTO question_tags (questionId, tag) VALUES ($1, $2) ON CONFLICT (questionId, tag) DO NOTHING",
		questionId, tag,
	)

	return err
}

func (r QuestionPostgresRepository) GetAllQuestionTags() ([]string, error) {
	rows, err := r.db.Query("SELECT DISTINCT(tag) FROM question_tags ORDER BY tag")
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

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

func (r QuestionPostgresRepository) GetTagsForQuestion(ID int) ([]string, error) {
	rows, err := r.db.Query("SELECT tag FROM question_tags WHERE questionId = $1", ID)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

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
