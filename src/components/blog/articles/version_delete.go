package articles

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
)

// DeleteVersionParams lists the params allowed by DeleteVersion
type DeleteVersionParams struct {
	ArticleID string `from:"url" json:"article_id" params:"uuid"`
	ID        string `from:"url" json:"id" params:"uuid"`
}

// DeleteVersion represents a API to soft delete a Version of an article
func DeleteVersion(req *router.Request) error {
	params := req.Params.(*DeleteVersionParams)

	a, err := GetByID(params.ArticleID)
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
	if v.ID == a.Version.ID {
		return httperr.NewConflict("cannot delete a version currently in use by an article")
	}

	if err := v.Delete(); err != nil {
		return err
	}

	req.NoContent()
	return nil
}
