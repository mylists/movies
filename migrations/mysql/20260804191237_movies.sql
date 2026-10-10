-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS movies (
    id CHAR(36) PRIMARY KEY,
    title VARCHAR(255) NOT NULL UNIQUE,
    release_year INTEGER DEFAULT 0,
    director VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS user_movies (
    user_id CHAR(36) NOT NULL,
    movie_id CHAR(36) NOT NULL,
    status VARCHAR(50) DEFAULT 'plan_to_watch',
    rating INTEGER DEFAULT 0,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, movie_id),
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE
);

CREATE INDEX idx_user_movies_user_id ON user_movies(user_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_movies;
DROP TABLE IF EXISTS movies;
-- +goose StatementEnd

