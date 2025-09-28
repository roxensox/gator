-- +goose Up
CREATE TABLE feed_follows (
	id UUID primary key,
	created_at TIMESTAMP,
	updated_at TIMESTAMP,
	u_id UUID,
	f_id UUID,
	FOREIGN KEY(u_id) REFERENCES users ON DELETE CASCADE,
	FOREIGN KEY(f_id) REFERENCES feeds ON DELETE CASCADE
);

-- +goose Down
DROP TABLE feed_follows;

