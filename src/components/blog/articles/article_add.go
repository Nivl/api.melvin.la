package articles

import (
	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/gosimple/slug"
)

// AddParams lists the params allowed by HandlerAdd
type AddParams struct {
	Title       string `from:"form" json:"title" params:"required,trim"`
	Subtitle    string `from:"form" json:"subtitle,omitempty"`
	Description string `from:"form" json:"description,omitempty"`
	Content     string `from:"form" json:"content,omitempty"`
}

// Add represents an API handler to add a new article
func Add(req *router.Request) error {
	params := req.Params.(*AddParams)

	tx, err := db.Writer.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create the article
	a := &Article{
		Slug:   slug.Make(params.Title),
		UserID: req.User.ID,
		User:   req.User,
	}
	if err := a.SaveQ(tx); err != nil {
		if db.IsDup(err) {
			return httperr.NewConflict("slug already in use")
		}
		return err
	}

	// Create the version
	v := &Version{
		Title:       params.Title,
		Subtitle:    params.Subtitle,
		Content:     params.Content,
		Description: params.Description,
		ArticleID:   a.ID,
	}
	if err := v.SaveQ(tx); err != nil {
		return err
	}

	// Set the Version in the article
	a.CurrentVersion = &v.ID
	a.Version = v
	if err := a.SaveQ(tx); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	req.Created(a.Export())
	return nil
}
