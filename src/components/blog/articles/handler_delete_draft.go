package articles

import (
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/router"
)

// HandlerDeleteDraftParams lists the params allowed by HandlerDelete
type HandlerDeleteDraftParams struct {
	ID string `from:"url" json:"id" params:"uuid"`
}

// HandlerDeleteDraft represents a API handler to soft delete a draft of an article
func HandlerDeleteDraft(req *router.Request) error {
	params := req.Params.(*HandlerDeleteDraftParams)

	a, err := Get(params.ID)
	if err != nil {
		return err
	}
	if a.IsZero() || req.User.ID != a.UserID {
		return apierror.NewNotFound()
	}

	if err := a.FetchDraft(); err != nil {
		return err
	}

	if a.Draft != nil {
		if err := a.Draft.Delete(); err != nil {
			return err
		}
	}

	req.NoContent()
	return nil
}
