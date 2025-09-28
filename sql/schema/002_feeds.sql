-- +goose Up
CREATE TABLE feeds (
	id UUID primary key,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	name VARCHAR UNIQUE NOT NULL,
	url VARCHAR UNIQUE NOT NULL,
	user_id UUID NOT NULL REFERENCES users ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feeds;
