package users

import (
	"fmt"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/formfile"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/storage/filestorage"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-rest-tools/types/filetype"
	"github.com/Nivl/go-rest-tools/types/ptrs"
)

var uploadPictureEndpoint = &router.Endpoint{
	Verb:    "PUT",
	Path:    "/users/{id}/picture",
	Handler: UploadPicture,
	Guard: &guard.Guard{
		Auth:        guard.LoggedUserAccess,
		ParamStruct: &UploadPictureParams{},
	},
}

// UploadPictureParams represents the params accepted by the UploadPicture endpoint
type UploadPictureParams struct {
	ID      string             `from:"url" json:"id" params:"required,uuid"`
	Picture *formfile.FormFile `from:"file" json:"picture" params:"required,image"`
}

// UploadPicture is an endpoint used to upload a picture
func UploadPicture(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*UploadPictureParams)
	currentUser := req.User()

	// Admin are allowed to update any users
	if !currentUser.IsAdmin && params.ID != currentUser.ID {
		return apierror.NewForbidden()
	}

	profile, err := GetByIDWithProfile(deps.DB, params.ID)
	if err != nil {
		return err
	}

	// we use the shasum as filename that way we don't re-upload the same
	// image twice
	shasum, err := filetype.SHA256Sum(params.Picture.File)
	if err != nil {
		return err
	}

	fileDest := fmt.Sprintf("users/pictures/%s", shasum)
	_, url, err := deps.Storage.WriteIfNotExist(params.Picture.File, fileDest)
	if err != nil {
		return err
	}

	// not a big deal if that fails
	deps.Storage.SetAttributes(fileDest, &filestorage.UpdatableFileAttributes{
		ContentType: params.Picture.Mime,
	})

	// Save the new picture
	profile.Picture = ptrs.NewString(url)
	if err = profile.Update(deps.DB); err != nil {
		return err
	}

	return req.Response().Ok(profile.ExportPrivate())
}
