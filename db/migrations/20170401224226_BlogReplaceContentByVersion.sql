
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
DROP TABLE blog_article_contents;

ALTER TABLE blog_articles ADD COLUMN current_version UUID REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;

CREATE TABLE blog_article_versions
(
  id UUID NOT NULL,
  article_id UUID NOT NULL REFERENCES blog_articles(id) ON DELETE CASCADE ON UPDATE CASCADE,

  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,

  title VARCHAR NOT NULL CHECK (length(title) < 255),
  subtitle VARCHAR NOT NULL DEFAULT '' CHECK (length(subtitle) < 255),
  description VARCHAR NOT NULL DEFAULT '' CHECK (length(subtitle) < 2000),
  content VARCHAR NOT NULL DEFAULT '' CHECK (length(subtitle) < 10000),

  PRIMARY KEY (id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

