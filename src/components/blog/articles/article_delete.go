package articles

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
)

// DeleteParams lists the params allowed by the Delete handler
type DeleteParams struct {
	ID string `from:"url" json:"id" params:"uuid"`
}

// Delete represents a API handler to soft delete a single article
func Delete(req *router.Request) error {
	params := req.Params.(*DeleteParams)

	a, err := GetByID(params.ID)
	if err != nil {
		return err
	}
	if a.IsZero() {
		return httperr.NewNotFound()
	}
	if a.PublishedAt != nil {
		return httperr.NewConflict("You cannot delete a published article")
	}

	if err := a.Delete(); err != nil {
		return err
	}

	req.NoContent()
	return nil
}
