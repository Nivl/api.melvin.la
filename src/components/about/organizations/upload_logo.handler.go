package organizations

import (
	"fmt"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/primitives/filetype"
	"github.com/Nivl/go-rest-tools/primitives/ptrs"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/formfile"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/storage/filestorage"
)

var uploadLogoEndpoint = &router.Endpoint{
	Verb:    "PUT",
	Path:    "/about/organizations/{id}/logo",
	Handler: UploadLogo,
	Guard: &guard.Guard{
		Auth:        guard.AdminAccess,
		ParamStruct: &UploadLogoParams{},
	},
}

// UploadLogoParams represents the params accepted by the UploadLogo endpoint
type UploadLogoParams struct {
	ID   string             `from:"url" json:"id" params:"required,uuid"`
	Logo *formfile.FormFile `from:"file" json:"logo" params:"required,image"`
}

// UploadLogo is an endpoint used to upload a logo
func UploadLogo(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*UploadLogoParams)

	org, err := GetByID(deps.DB, params.ID)
	if err != nil {
		return err
	}
	if org == nil {
		return httperr.NewNotFound()
	}

	// we use the shasum as filename that way we don't re-upload the same
	// image twice
	shasum, err := filetype.SHA256Sum(params.Logo.File)
	if err != nil {
		return err
	}

	fileDest := fmt.Sprintf("about/organizations/%s", shasum)
	_, url, err := deps.Storage.WriteIfNotExist(params.Logo.File, fileDest)
	if err != nil {
		return err
	}

	// not a big deal if that fails
	deps.Storage.SetAttributes(fileDest, &filestorage.UpdatableFileAttributes{
		ContentType: params.Logo.Mime,
	})

	// Save the new logo
	org.Logo = ptrs.NewString(url)
	if err = org.Update(deps.DB); err != nil {
		return err
	}

	return req.Response().Ok(org.Export())
}
