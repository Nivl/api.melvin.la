
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE user_profiles ADD COLUMN is_featured BOOLEAN DEFAULT NULL;
CREATE UNIQUE INDEX user_profiles_unique_is_featured_true ON user_profiles (is_featured)  WHERE (is_featured = true);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

