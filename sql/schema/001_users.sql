-- +goose Up
CREATE TABLE users(
	id UUID PRIMARY KEY,
	email VARCHAR(255) NOT NULL UNIQUE,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE users;
