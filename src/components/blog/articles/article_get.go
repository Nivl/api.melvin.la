package articles

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
)

// GetParams lists the params allowed by the Get Handler
type GetParams struct {
	ID string `from:"url" json:"id" params:"uuid"`
}

// Get represents a API handler to get a single article
func Get(req *router.Request) error {
	params := req.Params.(*GetParams)
	a, err := GetByID(params.ID)

	if err != nil {
		return err
	}
	if a.IsZero() {
		return httperr.NewNotFound()
	}
	// If it's not published, only the admin can access it, and we hide it from
	// the users
	if a.PublishedAt == nil && (req.User == nil || req.User.IsAdmin != true) {
		return httperr.NewNotFound()
	}

	req.Ok(a.Export())
	return nil
}
