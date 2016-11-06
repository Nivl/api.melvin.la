
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE users (
  uuid VARCHAR(36) NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,
  email VARCHAR(255) NOT NULL UNIQUE,
  name VARCHAR(255) NOT NULL,
  password VARCHAR(255) NOT NULL,
  PRIMARY KEY (uuid)
);

CREATE TABLE sessions (
  uuid VARCHAR(36) NOT NULL,
  user_uuid VARCHAR(36) NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
  created_at timestamptz NOT NULL,
  deleted_at timestamptz,
  PRIMARY KEY (uuid)
);

CREATE TABLE blog_articles (
  uuid VARCHAR(36) NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,
  user_uuid VARCHAR(36) NOT NULL REFERENCES users(uuid) ON DELETE CASCADE,
  title VARCHAR(255) NOT NULL,
  slug VARCHAR(255) NOT NULL UNIQUE,
  subtitle VARCHAR(255) NOT NULL  DEFAULT '',
  description TEXT NOT NULL  DEFAULT '',
  content TEXT NOT NULL DEFAULT '',
  is_published BOOLEAN NOT NULL DEFAULT false,
  PRIMARY KEY (uuid)
);

-- +goose Down

DROP TABLE blog_articles, sessions, users;

-- SQL section 'Down' is executed when this migration is rolled back

