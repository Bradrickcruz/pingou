-- +goose Up
ALTER TABLE incidents
ADD COLUMN duration_seconds INTEGER;

-- +goose Down
ALTER TABLE incidents
DROP COLUMN duration_seconds;