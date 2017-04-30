
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE blog_articles ALTER COLUMN slug TYPE VARCHAR;
ALTER TABLE blog_articles ADD CHECK (length(slug) < 260);

ALTER TABLE users ALTER COLUMN email TYPE VARCHAR;
ALTER TABLE users ADD CHECK (length(email) < 200);

ALTER TABLE users ALTER COLUMN name TYPE VARCHAR;
ALTER TABLE users ADD CHECK (length(name) < 200);

ALTER TABLE users ALTER COLUMN password TYPE VARCHAR;
ALTER TABLE users ADD CHECK (length(password) < 255);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

