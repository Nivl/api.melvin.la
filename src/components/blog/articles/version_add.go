package articles

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
)

// AddVersionParams lists the params allowed by AddVersion
type AddVersionParams struct {
	ArticleID string `from:"url" json:"article_id" params:"uuid"`
}

// AddVersion represents an API handler to add a new versions of an article
func AddVersion(req *router.Request) error {
	params := req.Params.(*AddVersionParams)

	a, err := GetByID(params.ArticleID)
	if err != nil {
		return err
	}
	if a.IsZero() {
		return httperr.NewNotFoundR("article not found")
	}

	// A new version is a dup of the current one
	v := a.Version
	v.ID = ""
	v.CreatedAt = nil
	v.UpdatedAt = nil

	if err := v.Create(); err != nil {
		return err
	}

	req.Ok(v.Export())
	return nil
}
