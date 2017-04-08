package articles

import (
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db"
)

// Article is a structure representing an article that can be saved in the database
//go:generate api-cli generate model Article -t blog_articles
type Article struct {
	ID             string   `db:"id"`
	Slug           string   `db:"slug"`
	CreatedAt      *db.Time `db:"created_at"`
	UpdatedAt      *db.Time `db:"updated_at"`
	DeletedAt      *db.Time `db:"deleted_at"`
	PublishedAt    *db.Time `db:"published_at"`
	UserID         string   `db:"user_id"`
	CurrentVersion *string  `db:"current_version"`

	*Version   `db:"version"`
	*auth.User `db:"user"`
}

// Articles represents a list of Articles
type Articles []Article
