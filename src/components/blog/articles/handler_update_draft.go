package articles

import (
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/db"
	"github.com/melvin-laplanche/ml-api/src/ptrs"
	"github.com/melvin-laplanche/ml-api/src/router"
)

// HandlerUpdateDraftParams lists the params allowed by HandlerGet
type HandlerUpdateDraftParams struct {
	ID string `from:"url" json:"id" params:"uuid"`

	Title       string `from:"form" json:"title" params:"trim"`
	Subtitle    string `from:"form" json:"subtitle" params:"trim"`
	Slug        string `from:"form" json:"slug" params:"trim"`
	Content     string `from:"form" json:"content" params:"trim"`
	Description string `from:"form" json:"description" params:"trim"`
	Promote     bool   `from:"form" json:"promote"`
}

// HandlerUpdateDraft is an API handler to update the draft of an article
func HandlerUpdateDraft(req *router.Request) {
	params := req.Params.(*HandlerUpdateDraftParams)

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
		req.Error(err)
		return
	}
	if a.IsZero() || req.User.ID != a.UserID {
		req.Error(apierror.NewNotFound())
		return
	}

	// fetch the current draft
	if err := a.FetchDraft(); err != nil {
		req.Error(err)
		return
	}

	draft := a.Draft

	if draft == nil {
		draft = a.Content.ToDraft()
		draft.ID = ""
		draft.IsDraft = ptrs.NewBool(true)
		draft.IsCurrent = ptrs.NewBool(false)
		a.Draft = draft
	}

	draftUpdated := false

	if params.Title != "" {
		draft.Title = params.Title
		draftUpdated = true
	}

	if params.Subtitle != "" {
		draft.Subtitle = params.Subtitle
		draftUpdated = true
	}

	if params.Content != "" {
		draft.Content = params.Content
		draftUpdated = true
	}

	if params.Description != "" {
		draft.Description = params.Description
		draftUpdated = true
	}

	tx, err := db.Con().Beginx()
	if err != nil {
		req.Error(err)
		return
	}

	if params.Promote {
		a.Content.IsCurrent = nil
		if err := a.Content.SaveTx(tx); err != nil {
			tx.Rollback()
			req.Error(err)
			return
		}

		draft.IsDraft = nil
		draft.IsCurrent = ptrs.NewBool(true)
		a.Content = draft.ToContent()
		a.Draft = nil
		draftUpdated = true
	}

	if draftUpdated {
		if err := draft.SaveTx(tx); err != nil {
			tx.Rollback()
			req.Error(err)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		req.Error(err)
		return
	}

	req.Ok(a.PrivateExport())
}
