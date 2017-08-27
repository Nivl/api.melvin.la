
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE user_profiles (
  id UUID NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,
  user_id uuid NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
  first_name VARCHAR CHECK (length(first_name) < 50),
  last_name VARCHAR CHECK (length(last_name) < 50),
  picture VARCHAR CHECK (length(picture) < 255),
  phone_number VARCHAR CHECK (length(phone_number) < 20),
  public_email VARCHAR CHECK (length(public_email) < 255),
  linkedin_custom_url VARCHAR CHECK (length(linkedin_custom_url) < 255),
  facebook_username VARCHAR CHECK (length(facebook_username) < 255),
  twitter_username VARCHAR CHECK (length(twitter_username) < 255),
  PRIMARY KEY (id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

