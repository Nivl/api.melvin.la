
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

DROP TABLE about_education_medias;
DROP TABLE about_education;

CREATE TABLE about_education (
  id UUID NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,

  organization_id UUID NOT NULL REFERENCES about_organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
  degree VARCHAR NOT NULL CHECK (length(degree) < 255),
  gpa VARCHAR CHECK (length(gpa) < 5),
  location VARCHAR CHECK (length(location) < 255),
  description VARCHAR CHECK (length(description) < 10000),
  start_year smallint NOT NULL CHECK (start_year > 1900 AND start_year < 2100),
  end_year smallint CHECK (end_year > 1900 AND end_year < 2100),
  PRIMARY KEY (id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

