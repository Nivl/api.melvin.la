package articles

import (
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/db"
	"github.com/melvin-laplanche/ml-api/src/router"
)

// HandlerGetParams lists the params allowed by HandlerGet
type HandlerGetParams struct {
	ID string `from:"url" json:"id" params:"uuid"`
}

// HandlerGet represents a API handler to get a single article
func HandlerGet(req *router.Request) error {
	params := req.Params.(*HandlerGetParams)

	a := &Article{}
	stmt := `SELECT articles.*,
                  ` + auth.UserJoinSQL("users") + `,
                  ` + ContentJoinSQL("content") + `
					FROM blog_articles articles
					JOIN users ON users.id = articles.user_id
					JOIN blog_article_contents content ON content.article_id = articles.id
					WHERE articles.deleted_at IS NULL
            AND articles.id = $1
						AND content.is_current IS true`

	if err := db.Get(a, stmt, params.ID); err != nil {
		return err
	}

	if a.IsZero() {
		return apierror.NewNotFound()
	}

	if req.User != nil && req.User.ID == a.UserID {
		// If the user is the author, let's get it's draft
		if err := a.FetchDraft(); err != nil {
			return err
		}

		req.Ok(a.PrivateExport())
		return nil
	}

	if a.PublishedAt == nil {
		return apierror.NewNotFound()
	}

	req.Ok(a.PublicExport())
	return nil
}
