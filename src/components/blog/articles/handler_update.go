package articles

import (
	"strconv"

	"github.com/gosimple/slug"
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/db"
	"github.com/melvin-laplanche/ml-api/src/ptrs"
	"github.com/melvin-laplanche/ml-api/src/router"
)

// HandlerUpdateParams lists the params allowed by HandlerGet
type HandlerUpdateParams struct {
	ID string `from:"url" json:"id" params:"uuid"`

	Title       string `from:"form" json:"title" params:"trim"`
	Subtitle    string `from:"form" json:"subtitle" params:"trim"`
	Slug        string `from:"form" json:"slug" params:"trim"`
	Content     string `from:"form" json:"content" params:"trim"`
	Description string `from:"form" json:"description" params:"trim"`
	Publish     string `from:"form" json:"publish" params:"bool"`
}

// HandlerUpdate represents a API handler to update a single article
func HandlerUpdate(req *router.Request) error {
	params := req.Params.(*HandlerUpdateParams)

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
	if a.IsZero() || req.User.ID != a.UserID {
		return apierror.NewNotFound()
	}

	articleUpdated := false
	contentUpdated := false

	// We create a new content for versioning
	newContent := *a.Content
	newContent.ID = ""

	if params.Title != "" {
		newContent.Title = params.Title
		contentUpdated = true
	}

	if params.Subtitle != "" {
		newContent.Subtitle = params.Subtitle
		contentUpdated = true
	}

	if params.Content != "" {
		newContent.Content = params.Content
		contentUpdated = true
	}

	if params.Description != "" {
		newContent.Description = params.Description
		contentUpdated = true
	}

	if params.Publish != "" {
		publish, _ := strconv.ParseBool(params.Publish)

		if publish && a.PublishedAt == nil {
			a.PublishedAt = db.Now()
		} else if !publish && a.PublishedAt != nil {
			a.PublishedAt = nil
		}

		articleUpdated = true
	}

	if params.Slug != "" {
		a.Slug = slug.Make(params.Slug)

		if a.Slug != params.Slug {
			return apierror.NewBadRequest("invalid value for slug: %s did you mean %s", params.Slug, a.Slug)
		}

		articleUpdated = true
	}

	tx, err := db.Con().Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if contentUpdated {
		a.Content.IsCurrent = ptrs.NewBool(false)
		if err := a.Content.UpdateTx(tx); err != nil {
			return err
		}

		newContent.IsCurrent = ptrs.NewBool(true)
		if err := newContent.CreateTx(tx); err != nil {
			return err
		}
		a.Content = &newContent
	}

	if articleUpdated {
		if err := a.UpdateTx(tx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	req.Ok(a.PrivateExport())
	return nil
}
