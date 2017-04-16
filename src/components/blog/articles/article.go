package articles

import (
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db"
)

// Article is a structure representing an article that can be saved in the database
//go:generate api-cli generate model Article -t blog_articles -e Get
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
	*auth.User `db:"users"`
}

// Articles represents a list of Articles
type Articles []Article

// GetByID finds and returns an active article by ID
func GetByID(id string) (*Article, error) {
	a := &Article{}
	stmt := `SELECT articles.*,
						` + auth.UserJoinSQL("users") + `,
						` + JoinVersionSQL("version") + `
						FROM blog_articles articles
						JOIN users ON users.id = articles.user_id
						JOIN blog_article_versions version ON version.id = articles.current_version
						WHERE articles.id=$1
						AND articles.deleted_at IS NULL
						LIMIT 1`
	err := db.Get(a, stmt, id)
	// We want to return nil if a article is not found
	if a.ID == "" {
		return nil, err
	}
	return a, err
}
