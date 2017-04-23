
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE INDEX blog_article_versions_created_at_idx ON blog_article_versions (created_at);
CREATE INDEX blog_article_versions_deleted_at_idx ON blog_article_versions (deleted_at);
CREATE INDEX blog_article_versions_article_id_idx ON blog_article_versions (article_id);

CREATE INDEX blog_articles_created_at_idx ON blog_articles (created_at);
CREATE INDEX blog_articles_current_version_idx ON blog_articles (current_version);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

