
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TYPE media_type AS ENUM ('image', 'video', 'link');

CREATE TABLE about_organizations (
  id UUID NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,

  name VARCHAR NOT NULL UNIQUE CHECK (length(name) < 255),
  short_name VARCHAR UNIQUE CHECK (length(short_name) < 255),
  logo VARCHAR(255) CHECK (length(logo) < 255),
  website VARCHAR(255) CHECK (length(website) < 255),
  PRIMARY KEY (id)
);


CREATE TABLE about_experience (
  id UUID NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,

  organization_id UUID NOT NULL REFERENCES about_organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
  job_title VARCHAR NOT NULL CHECK (length(job_title) < 255),
  location VARCHAR CHECK (length(location) < 255),
  description VARCHAR CHECK (length(description) < 10000),
  start_date DATE NOT NULL,
  end_date DATE,
  PRIMARY KEY (id)
);

CREATE TABLE about_experience_medias (
  id UUID NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,

  experience_id UUID NOT NULL REFERENCES about_experience(id) ON DELETE CASCADE ON UPDATE CASCADE,
  position SMALLINT NOT NULL UNIQUE,
  type media_type NOT NULL,
  content VARCHAR NOT NULL CHECK (length(content) < 255),
  image_preview VARCHAR NOT NULL CHECK (length(image_preview) < 255),
  PRIMARY KEY (id)
);

CREATE TABLE about_education (
  id UUID NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,

  organization_id UUID NOT NULL REFERENCES about_organizations(id) ON DELETE CASCADE ON UPDATE CASCADE,
  degree VARCHAR NOT NULL CHECK (length(degree) < 255),
  gpa VARCHAR NOT NULL CHECK (length(gpa) < 5),
  location VARCHAR CHECK (length(location) < 255),
  description VARCHAR CHECK (length(description) < 10000),
  start_year smallint NOT NULL CHECK (start_year > 1900 AND start_year < 2100),
  end_year smallint CHECK (end_year > 1900 AND end_year < 2100),
  PRIMARY KEY (id)
);

CREATE TABLE about_education_medias (
  id UUID NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,

  education_id UUID NOT NULL REFERENCES about_education(id) ON DELETE CASCADE ON UPDATE CASCADE,
  position SMALLINT NOT NULL UNIQUE,
  type media_type NOT NULL,
  content VARCHAR NOT NULL CHECK (length(content) < 255),
  image_preview VARCHAR NOT NULL CHECK (length(image_preview) < 255),
  PRIMARY KEY (id)
);

CREATE TABLE about_tech (
  id UUID NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,

  logo VARCHAR NOT NULL CHECK (length(logo) < 255),
  position SMALLINT NOT NULL UNIQUE,
  PRIMARY KEY (id)
);

CREATE TABLE about_movies (
  id UUID NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,

  poster VARCHAR NOT NULL CHECK (length(poster) < 255),
  name VARCHAR NOT NULL CHECK (length(name) < 255),
  original_name VARCHAR NOT NULL CHECK (length(original_name) < 255),
  imdb_link VARCHAR NOT NULL CHECK (length(imdb_link) < 255),
  rotten_tomatoes_link VARCHAR NOT NULL CHECK (length(rotten_tomatoes_link) < 255),
  position SMALLINT NOT NULL UNIQUE,
  PRIMARY KEY (id)
);

CREATE TABLE about_musics (
  id UUID NOT NULL,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz,

  poster VARCHAR NOT NULL CHECK (length(poster) < 255),
  name VARCHAR NOT NULL CHECK (length(name) < 255),
  artist VARCHAR NOT NULL CHECK (length(artist) < 255),
  album VARCHAR NOT NULL CHECK (length(album) < 255),
  spotify_link VARCHAR CHECK (length(spotify_link) < 255),
  google_music_link VARCHAR NOT NULL CHECK (length(google_music_link) < 255),
  youtube_link VARCHAR NOT NULL CHECK (length(youtube_link) < 255),
  position SMALLINT NOT NULL UNIQUE,
  PRIMARY KEY (id)
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

