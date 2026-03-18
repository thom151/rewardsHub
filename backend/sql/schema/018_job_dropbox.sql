-- +goose Up
ALTER TABLE job
ADD COLUMN dropbox_folder_path TEXT;

-- +goose Down
ALTER TABLE job
DROP COLUMN dropbox_folder_path;

