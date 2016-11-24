
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE sessions ADD COLUMN updated_at timestamptz;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

