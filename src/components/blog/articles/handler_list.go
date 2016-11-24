package articles

import (
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/router"
)

// HandlerList represents a API handler to get a list of articles
func HandlerList(req *router.Request) {
	arts := Articles{}

	// We want the published articles and the unpublished that we own
	var args []interface{}
	visibilityStmt := "articles.is_published IS TRUE"
	if req.User != nil {
		visibilityStmt += ` OR articles.user_id = $1`
		args = append(args, req.User.ID)
	}

	stmt := `SELECT articles.*, ` + auth.UserJoinSQL("users") + `
					FROM blog_articles articles
					LEFT JOIN users ON users.id = articles.user_id
					WHERE articles.deleted_at IS NULL
						AND ( ` + visibilityStmt + ` )
					ORDER BY articles.created_at`
	if err := sql().Select(&arts, stmt, args...); err != nil {
		req.Error(err)
		return
	}

	req.Ok(arts.Export())
}
