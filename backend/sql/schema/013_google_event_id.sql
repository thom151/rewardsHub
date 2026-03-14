-- +goose Up
ALTER TABLE job
ADD COLUMN google_event_id TEXT;

-- +goose Down
ALTER TABLE job
DROP COLUMN google_event_id;
