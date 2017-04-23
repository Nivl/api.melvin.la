package articles

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/storage/db"
)

// ListVersionParams lists the params allowed by UpdateVersion
type ListVersionParams struct {
	ArticleID string `from:"url" json:"article_id" params:"uuid"`
}

// ListVersion represents an API handler to list all the versions of an article
func ListVersion(req *router.Request) error {
	params := req.Params.(*ListVersionParams)

	a, err := GetByID(params.ArticleID)
	if err != nil {
		return err
	}
	if a.IsZero() {
		return httperr.NewNotFoundR("article not found")
	}

	var versions Versions
	stmt := "SELECT * FROM blog_article_versions WHERE article_id=$1 ORDER BY created_at ASC"
	if err := db.Writer.Select(&versions, stmt, a.ID); err != nil {
		return err
	}
	req.Ok(versions.Export())
	return nil
}
