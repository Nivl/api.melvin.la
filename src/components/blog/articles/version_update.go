package articles

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
)

// UpdateVersionParams lists the params allowed by UpdateVersion
type UpdateVersionParams struct {
	ID        string `from:"url" json:"id" params:"uuid"`
	ArticleID string `from:"url" json:"article_id" params:"uuid"`

	Title       string `from:"form" json:"title" params:"trim"`
	Subtitle    string `from:"form" json:"subtitle" params:"trim"`
	Content     string `from:"form" json:"content" params:"trim"`
	Description string `from:"form" json:"description" params:"trim"`
}

// UpdateVersion is an handler to update a version
func UpdateVersion(req *router.Request) error {
	params := req.Params.(*UpdateVersionParams)

	v, err := GetVersionByID(params.ID)
	if err != nil {
		return err
	}
	if v.IsZero() {
		return httperr.NewNotFound()
	}
	if v.ArticleID != params.ArticleID {
		return httperr.NewConflict("article %s doesn't match with provided version %s",
			params.ArticleID, params.ID)
	}

	if params.Title != "" {
		v.Title = params.Title
	}
	if params.Subtitle != "" {
		v.Subtitle = params.Subtitle
	}
	if params.Content != "" {
		v.Content = params.Content
	}
	if params.Description != "" {
		v.Description = params.Description
	}

	if err := v.Save(); err != nil {
		return err
	}

	req.Ok(v.Export())
	return nil
}
