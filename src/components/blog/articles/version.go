package articles

import (
	"github.com/melvin-laplanche/ml-api/src/db"
)

// Version is a structure representing a version of an article
//go:generate api-cli generate model Version -t blog_article_versions
type Version struct {
	ID        string   `db:"id"`
	CreatedAt *db.Time `db:"created_at"`
	UpdatedAt *db.Time `db:"updated_at"`
	DeletedAt *db.Time `db:"deleted_at"`

	ArticleID string `db:"article_id"`

	Title       string `db:"title"`
	Content     string `db:"content"`
	Subtitle    string `db:"subtitle"`
	Description string `db:"description"`
}
