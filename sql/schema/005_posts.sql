-- +goose Up
CREATE TABLE posts (
	id UUID PRIMARY KEY,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	title VARCHAR,
	url VARCHAR UNIQUE,
	description VARCHAR,
	published_at TIMESTAMP,
	feed_id UUID NOT NULL REFERENCES feeds ON DELETE CASCADE
);

-- +goose Down
DROP TABLE posts;
