package articles

import (
	"fmt"
	"testing"

	"github.com/dchest/uniuri"
	"github.com/gosimple/slug"
	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/db"
)

// Article is a structure representing an article that can be saved in the database
//go:generate api-cli generate model Article -t blog_articles -e Create,Update
type Article struct {
	ID          string   `db:"id"`
	Title       string   `db:"title"`
	Content     string   `db:"content"`
	Slug        string   `db:"slug"`
	Subtitle    string   `db:"subtitle"`
	Description string   `db:"description"`
	CreatedAt   *db.Time `db:"created_at"`
	UpdatedAt   *db.Time `db:"updated_at"`
	DeletedAt   *db.Time `db:"deleted_at"`
	IsPublished bool     `db:"is_published"`
	UserID      string   `db:"user_id"`
	auth.User   `db:"users"`
}

// Articles represents a list of Articles
type Articles []Article

// Create persists an article in the database
func (a *Article) Create() error {
	if a == nil {
		return apierror.NewServerError("article not instanced")
	}

	if a.Slug == "" {
		a.Slug = slug.Make(a.Title)
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

// NewTestArticle returns a published article with random values
func NewTestArticle(t *testing.T, a *Article) (*Article, *auth.User, *auth.Session) {
	if a == nil {
		a = &Article{
			IsPublished: true,
		}
	}

	if a.Title == "" {
		a.Title = uniuri.New()
	}

	var user *auth.User
	var session *auth.Session

	if a.UserID == "" {
		user, session = auth.NewTestAuth(t)
		a.User = *user
		a.UserID = user.ID
	}

	if err := a.Create(); err != nil {
		t.Fatalf("failed to save article: %s", err)
	}
	return a, user, session
}
