-- +goose Up
CREATE TABLE feed_follows (
	id UUID primary key,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	u_id UUID NOT NULL,
	f_id UUID NOT NULL,
	FOREIGN KEY(u_id) REFERENCES users ON DELETE CASCADE,
	FOREIGN KEY(f_id) REFERENCES feeds ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feed_follows;

