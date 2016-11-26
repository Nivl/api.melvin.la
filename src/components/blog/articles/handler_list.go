package articles

import (
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/router"
)

// HandlerListParams lists the params allowed by HandlerList
type HandlerListParams struct {
	// 2 values possible all, own
	Target string `from:"query" json:"target,omitempty" params:"trim" default:"all"`
}

// HandlerList represents a API handler to get a list of articles
func HandlerList(req *router.Request) {
	params := req.Params.(*HandlerListParams)
	var output *Payloads
	var err error

	switch params.Target {
	case "all":
		output, err = listAll()
	// case "own":
	// 	output, err = listOwn(req.User)
	default:
		err = apierror.NewBadRequest(`unknown value "%s" for param "target"`, params.Target)
	}

	if err != nil {
		req.Error(err)
	}

	req.Ok(output)
}

func listAll() (*Payloads, error) {
	arts := Articles{}

	stmt := `SELECT articles.*,
                  ` + auth.UserJoinSQL("users") + `,
                  ` + ContentJoinSQL("content") + `
					FROM blog_articles articles
					JOIN users ON users.id = articles.user_id
					JOIN blog_article_contents content ON content.article_id = articles.id
					WHERE articles.deleted_at IS NULL
						AND articles.published_at IS NOT NULL
						AND content.is_current IS TRUE
					ORDER BY articles.created_at`
	if err := sql().Select(&arts, stmt); err != nil {
		return nil, err
	}

	return arts.PublicExport(), nil
}

// func listOwn(user *auth.User) (*Payloads, error) {

// 	arts := Articles{}

// 	// We want the published articles and the unpublished that we own
// 	var args []interface{}
// 	visibilityStmt := "articles.is_published IS TRUE"
// 	if req.User != nil {
// 		visibilityStmt += ` OR articles.user_id = $1`
// 		args = append(args, req.User.ID)
// 	}

// 	stmt := `SELECT articles.*,
//                   ` + auth.UserJoinSQL("user") + `
//                   ` + auth.UserJoinSQL("content") + `
// 					FROM blog_articles articles
// 					LEFT JOIN users ON user.id = articles.user_id
// 					WHERE articles.deleted_at IS NULL
// 						AND ( ` + visibilityStmt + ` )
// 					ORDER BY articles.created_at`
// 	if err := sql().Select(&arts, stmt, args...); err != nil {
// 		req.Error(err)
// 		return
// 	}

// 	return arts.Export()
// }
