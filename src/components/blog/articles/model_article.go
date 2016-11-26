package articles

import (
	"fmt"

	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/db"
)

// Article is a structure representing an article that can be saved in the database
//go:generate api-cli generate model Article -t blog_articles -e Create,Update
type Article struct {
	ID          string   `db:"id"`
	Slug        string   `db:"slug"`
	CreatedAt   *db.Time `db:"created_at"`
	UpdatedAt   *db.Time `db:"updated_at"`
	DeletedAt   *db.Time `db:"deleted_at"`
	PublishedAt *db.Time `db:"published_at"`
	UserID      string   `db:"user_id"`

	*auth.User `db:"users"`
	*Content   `db:"content"`
	*Draft     `db:"draft"`
}

// Articles represents a list of Articles
type Articles []Article

// Create persists an article in the database
func (a *Article) Create() error {
	if a == nil {
		return apierror.NewServerError("article not instanced")
	}

	if a.Slug == "" {
		return apierror.NewServerError("cannot persist an article with no slug")
	}

	// To prevent duplicates on the slug, we'll retry the insert() up to 10 times
	originalSlug := a.Slug
	var err error
	for i := 0; i < 10; i++ {
		err = a.doCreate()

		if err != nil {
			if db.SQLIsDup(err) == false {
				return apierror.NewServerError(err.Error())
			}

			// In case of duplicate we'll add "-X" at the end of the slug, where X is
			// a number
			a.Slug = fmt.Sprintf("%s-%d", originalSlug, i)
		} else {
			// everything went well
			return nil
		}
	}

	// after 10 try we just return an error
	return apierror.NewConflict(err.Error())
}

// Update updates most of the fields of a persisted user.
// Excluded fields are id, created_at, deleted_at
func (a *Article) Update() error {
	return nil
}
