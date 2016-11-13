package articles

import (
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/router"
)

// HandlerList represents a API handler to get a list of articles
func HandlerList(req *router.Request) {
	arts := Articles{}

	stmt := `SELECT articles.*, ` + auth.UserForeignSelect("users") + `
					FROM blog_articles articles
					LEFT JOIN users ON users.id = articles.user_id
					WHERE articles.deleted_at IS NULL
					ORDER BY articles.created_at`
	if err := sql().Select(&arts, stmt); err != nil {
		req.Error(err)
		return
	}

	req.Ok(arts.Export())
}
