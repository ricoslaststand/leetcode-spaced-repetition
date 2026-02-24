-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS problems (
    id INTEGER PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL UNIQUE,
    description TEXT,
    difficulty TEXT NOT NULL CHECK (difficulty IN ('easy', 'medium', 'hard')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS problem_topics (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    problem_id INTEGER NOT NULL,
    topic VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (problem_id, topic),
    FOREIGN KEY (problem_id) REFERENCES problems(id)
);

CREATE TABLE IF NOT EXISTS problem_submissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    problem_id INT NOT NULL,
    submission_date DATE NOT NULL,
    confidence_level TEXT NOT NULL CHECK (confidence_level IN ('again', 'hard', 'good', 'easy')),
    time_taken INTERVAL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (user_id, problem_id, submission_date),
    FOREIGN KEY (problem_id) REFERENCES problems(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS problem_submissions;
DROP TABLE IF EXISTS problem_topics;
DROP TABLE IF EXISTS problems;
-- +goose StatementEnd
