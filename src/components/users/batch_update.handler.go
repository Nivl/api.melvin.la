package users

import (
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-types/ptrs"
)

var batchUpdateEndpoint = &router.Endpoint{
	Verb:    "PATCH",
	Path:    "/users",
	Handler: Update,
	Guard: &guard.Guard{
		ParamStruct: &BatchUpdateParams{},
		Auth:        guard.AdminAccess,
	},
}

// BatchUpdateParams represents the params accepted Update to update a user
type BatchUpdateParams struct {
	FeaturedUser *string `from:"form" json:"featured_user"  params:"uuid,noempty"`
}

// BatchUpdate is a HTTP handler used to update a user
func BatchUpdate(req router.HTTPRequest, deps *router.Dependencies) error {
	params := req.Params().(*BatchUpdateParams)
	updatedProfiles := map[string]*Profile{}

	if params.FeaturedUser != nil {
		err := batchUpdateFeaturedUser(updatedProfiles, deps.DB, *params.FeaturedUser)
		if err != nil {
			return err
		}
	}

	profiles := Profiles{}
	for _, p := range updatedProfiles {
		profiles = append(profiles, p)
	}

	return req.Response().Ok(profiles.ExportPrivate())
}

func batchUpdateFeaturedUser(updatedUsers map[string]*Profile, dbCon db.Connection, userUUID string) error {
	// Retreive the user and the attached profile
	newFeatured, err := GetByIDWithProfile(dbCon, userUUID)
	if err != nil {
		return err
	}

	// We make sure the user is not already featured
	if newFeatured.IsFeatured != nil && *newFeatured.IsFeatured {
		return apierror.NewConflictR("featured_user", "user already featured")
	}

	currentFeatured, err := GetFeaturedProfile(dbCon)
	if err != nil && !apierror.IsNotFound(err) {
		return err
	}

	// If there're no user currently featured, we just use set the new one
	if apierror.IsNotFound(err) {
		newFeatured.IsFeatured = ptrs.NewBool(true)
		if err := newFeatured.Save(dbCon); err != nil {
			return err
		}
		updatedUsers[newFeatured.ID] = newFeatured
		return nil
	}

	// Create a transaction to make sure we don't endup without featured
	// users in case of database issue
	tx, err := dbCon.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Downgrade the current featured user
	currentFeatured.IsFeatured = nil
	if err := currentFeatured.Save(tx); err != nil {
		return err
	}
	// Upgrade the new featured user
	newFeatured.IsFeatured = ptrs.NewBool(true)
	if err := newFeatured.Save(tx); err != nil {
		return err
	}

	// Persist the changes
	if err := tx.Commit(); err != nil {
		return err
	}

	// cache the changes
	updatedUsers[currentFeatured.ID] = currentFeatured
	updatedUsers[newFeatured.ID] = newFeatured
	return nil
}
