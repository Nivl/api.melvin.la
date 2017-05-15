package articles

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/gosimple/slug"
)

// UpdateParams lists the params allowed by Update
type UpdateParams struct {
	ID      string `from:"url" json:"id" params:"uuid"`
	Slug    string `from:"form" json:"slug" params:"trim"`
	Version string `from:"form" json:"version" params:"uuid"`
	Publish *bool  `from:"form" json:"publish"`
}

// Update represents a API handler to update a single article
func Update(req *router.Request) error {
	params := req.Params.(*UpdateParams)

	a, err := GetByID(params.ID)
	if err != nil {
		return err
	}
	if a.IsZero() {
		return httperr.NewNotFoundR("article not found")
	}

	// Update the new version
	if params.Version != "" && params.Version != *a.CurrentVersion {
		v, err := GetVersionByID(params.Version)
		if err != nil {
			return err
		}
		if v == nil || v.ArticleID != a.ID {
			return httperr.NewNotFoundR("version not found")
		}
		a.CurrentVersion = &v.ID
		a.Version = v
	}

	// Publish / Un-publish
	if params.Publish != nil {
		if *params.Publish && a.PublishedAt == nil {
			a.PublishedAt = db.Now()
		} else if !*params.Publish && a.PublishedAt != nil {
			a.PublishedAt = nil
		}
	}

	// Update Slug
	if params.Slug != "" {
		a.Slug = slug.Make(params.Slug)
		if a.Slug != params.Slug {
			return httperr.NewBadRequest("invalid value for slug: %s did you mean %s", params.Slug, a.Slug)
		}
	}

	if err := a.Update(); err != nil {
		return err
	}

	req.Ok(a.Export())
	return nil
}
