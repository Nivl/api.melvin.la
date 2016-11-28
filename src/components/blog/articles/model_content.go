package articles

import (
	"fmt"

	"github.com/melvin-laplanche/ml-api/src/db"
)

// Content is a structure representing an article content that can be saved in the database
//go:generate api-cli generate model Content -t blog_article_contents -e JoinSQL
type Content struct {
	ID        string   `db:"id"`
	CreatedAt *db.Time `db:"created_at"`
	UpdatedAt *db.Time `db:"updated_at"`
	DeletedAt *db.Time `db:"deleted_at"`

	ArticleID string `db:"article_id"`
	IsCurrent *bool  `db:"is_current"`
	IsDraft   *bool  `db:"is_draft"`

	Title       string `db:"title"`
	Content     string `db:"content"`
	Subtitle    string `db:"subtitle"`
	Description string `db:"description"`
}

// ContentJoinSQL returns a string ready to be embed in a JOIN query
func ContentJoinSQL(prefix string) string {
	fields := []string{"id", "created_at", "updated_at", "deleted_at", "article_id", "is_current", "is_draft", "title", "content", "subtitle", "description"}
	output := ""

	for i, field := range fields {
		if i != 0 {
			output += ", "
		}

		fullName := fmt.Sprintf("%s.%s", prefix, field)
		output += fmt.Sprintf("%s \"%s\"", fullName, fullName)
	}

	return output
}

// ToDraft returns a Draft from a Content
func (c *Content) ToDraft() *Draft {
	draft := Draft(*c)
	return &draft
}
