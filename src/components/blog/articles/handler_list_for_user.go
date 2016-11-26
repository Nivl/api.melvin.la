package articles

import (
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/router"
)

// HandlerListForUserParams lists the params allowed by HandlerList
type HandlerListForUserParams struct {
	UserID string `from:"url" json:"user_id" params:"uuid"`
}

// HandlerListForUser represents a API handler to get a list of articles
func HandlerListForUser(req *router.Request) {
	params := req.Params.(*HandlerListForUserParams)

	user, err := auth.GetUser(params.UserID)
	if err != nil {
		req.Error(err)
		return
	}
	if user == nil {
		req.Error(apierror.NewNotFoundR("no such user"))
		return
	}

	var arts Articles
	stmt := `SELECT articles.*,
                  ` + auth.UserJoinSQL("users") + `,
                  ` + ContentJoinSQL("content") + `
					FROM blog_articles articles
					JOIN users ON users.id = articles.user_id
					JOIN blog_article_contents content ON content.article_id = articles.id
					WHERE articles.deleted_at IS NULL
						AND articles.published_at IS NOT NULL
						AND content.is_current IS TRUE
            AND users.id = $1
					ORDER BY articles.created_at`

	if err := sql().Select(&arts, stmt, params.UserID); err != nil {
		req.Error(err)
		return
	}

	req.Ok(arts.PublicExport())
}
