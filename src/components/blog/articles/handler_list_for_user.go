package articles

import (
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/db"
	"github.com/melvin-laplanche/ml-api/src/router"
)

// HandlerListForUserParams lists the params allowed by HandlerList
type HandlerListForUserParams struct {
	UserID string `from:"url" json:"user_id" params:"uuid"`
}

// HandlerListForUser represents a API handler to get a list of articles
func HandlerListForUser(req *router.Request) error {
	params := req.Params.(*HandlerListForUserParams)

	user, err := auth.GetUser(params.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return apierror.NewNotFoundR("no such user")
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

	if err := db.Con().Select(&arts, stmt, params.UserID); err != nil {
		return err
	}

	req.Ok(arts.PublicExport())
	return nil
}
