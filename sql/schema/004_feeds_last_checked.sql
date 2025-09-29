-- +goose Up
ALTER TABLE feeds ADD COLUMN last_checked TIMESTAMP;

-- +goose Down
ALTER TABLE feeds DROP last_checked;
