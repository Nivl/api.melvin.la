package articles

import (
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/router"
)

// HandlerDeleteParams lists the params allowed by HandlerDelete
type HandlerDeleteParams struct {
	ID string `from:"url" json:"id" params:"uuid"`
}

// HandlerDelete represents a API handler to soft delete a single article
func HandlerDelete(req *router.Request) error {
	params := req.Params.(*HandlerDeleteParams)

	a, err := Get(params.ID)
	if err != nil {
		return err
	}
	if a.IsZero() || req.User.ID != a.UserID {
		return apierror.NewNotFound()
	}

	if err := a.Delete(); err != nil {
		return err
	}

	req.NoContent()
	return nil
}
