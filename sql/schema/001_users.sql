-- +goose Up
CREATE TABLE users (
	id UUID primary key,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	name VARCHAR UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE users;
