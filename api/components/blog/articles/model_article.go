package articles

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/auth"
	"github.com/Nivl/api.melvin.la/api/db"
	"github.com/dchest/uniuri"
	"github.com/gosimple/slug"
	uuid "github.com/satori/go.uuid"
)

// Article is a structure representing an article that can be saved in the database
type Article struct {
	ID          string     `db:"id"`
	Title       string     `db:"title"`
	Content     string     `db:"content"`
	Slug        string     `db:"slug"`
	Subtitle    string     `db:"subtitle"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
	IsPublished bool       `db:"is_published"`
	UserID      string     `db:"user_id"`
	auth.User   `db:"users"`
}

// FullyDelete removes an article from the database
func (a *Article) FullyDelete() error {
	if a == nil {
		return errors.New("article not instanced")
	}

	if a.ID == "" {
		return errors.New("article has not been saved")
	}

	_, err := sql().Exec("DELETE FROM blog_articles WHERE id=$1", a.ID)
	return err
}

// Save creates or updates the article depending on the value of the id
func (a *Article) Save() error {
	if a == nil {
		return errors.New("article not instanced")
	}

	if a.ID == "" {
		return a.Create()
	}

	return a.Update()
}

// Create persists an article in the database
func (a *Article) Create() error {
	if a == nil {
		return apierror.NewServerError("article not instanced")
	}

	if a.Slug == "" {
		a.Slug = slug.Make(a.Title)
	}

	a.CreatedAt = time.Now()

	// To prevent duplicates on the slug, we'll retry the insert() up to 10 times
	originalSlug := a.Slug
	var err error
	for i := 0; i < 10; i++ {
		a.ID = uuid.NewV4().String()

		stmt := `INSERT INTO blog_articles
		(id, created_at, updated_at, title, content, slug, subtitle, description, is_published, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
		_, err := sql().Exec(stmt, a.ID, a.CreatedAt, a.UpdatedAt, a.Title, a.Content, a.Slug, a.Subtitle, a.Description, a.IsPublished, a.User.ID)

		if err != nil {
			// In case of duplicate we'll add "-X" at the end of the slug, where X is
			// a number
			a.Slug = fmt.Sprintf("%s-%d", originalSlug, i)

			if db.SQLIsDup(err) == false {
				return apierror.NewServerError(err.Error())
			}
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

	user, session := auth.NewTestAuth(t)
	a.User = *user

	if err := a.Save(); err != nil {
		t.Fatalf("failed to save article: %s", err)
	}
	return a, user, session
}
