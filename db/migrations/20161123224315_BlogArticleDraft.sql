
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

ALTER TABLE blog_articles ADD COLUMN published_at timestamptz;

ALTER TABLE blog_articles DROP COLUMN is_published;
ALTER TABLE blog_articles DROP COLUMN title;
ALTER TABLE blog_articles DROP COLUMN subtitle;
ALTER TABLE blog_articles DROP COLUMN description;
ALTER TABLE blog_articles DROP COLUMN content;

CREATE TABLE blog_article_contents (
  id UUID NOT NULL,

  article_id UUID NOT NULL REFERENCES blog_articles(id) ON DELETE CASCADE ON UPDATE CASCADE,
  is_current bool,
  is_draft bool,

  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,


  title VARCHAR(255) NOT NULL,
  subtitle VARCHAR(255) NOT NULL DEFAULT '',
  description TEXT NOT NULL DEFAULT '',
  content TEXT NOT NULL DEFAULT '',

  PRIMARY KEY (id),
  UNIQUE (article_id, is_current),
  UNIQUE (article_id, is_draft)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

