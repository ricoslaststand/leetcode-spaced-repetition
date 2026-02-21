-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS questions (
    id INTEGER PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    difficulty TEXT NOT NULL CHECK (difficulty IN ('easy', 'mild', 'hard')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS question_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    question_id INTEGER NOT NULL,
    tag VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (question_id, tag),
    FOREIGN KEY (question_id) REFERENCES questions(id)
);

CREATE TABLE IF NOT EXISTS question_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    question_id INT NOT NULL,
    submission_date DATE NOT NULL,
    confidence_level TEXT NOT NULL CHECK (confidence_level IN ('again', 'hard', 'good', 'easy')),
    time_taken INTERVAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, question_id, submission_date),
    FOREIGN KEY (question_id) REFERENCES questions(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS question_submissions;
DROP TABLE IF EXISTS question_tags;
DROP TABLE IF EXISTS questions;
-- +goose StatementEnd
