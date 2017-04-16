package articletestdata

import (
	"testing"

	"github.com/Nivl/go-rest-tools/primitives/models/lifecycle"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/security/auth/testdata"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/dchest/uniuri"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles"
	uuid "github.com/satori/go.uuid"
)

// NewArticle returns a published article with random values
func NewArticle(t *testing.T, a *articles.Article) (*articles.Article, *auth.User, *auth.Session) {
	if a == nil {
		a = &articles.Article{
			PublishedAt: db.Now(),
		}
	}

	if a.Slug == "" {
		a.Slug = uuid.NewV4().String()
	}

	var user *auth.User
	var session *auth.Session

	if a.UserID == "" {
		user, session = testdata.NewAuth(t)
		user.IsAdmin = true
		user.Update()
		a.User = user
		a.UserID = user.ID
	}

	if err := a.Create(); err != nil {
		t.Fatalf("failed to save article: %s", err)
	}

	a.Version = NewVersion(t, a, a.Version)
	a.CurrentVersion = &a.Version.ID
	if err := a.Update(); err != nil {
		t.Fatalf("failed to update article: %s", err)
	}

	lifecycle.SaveModels(t, a)
	return a, user, session
}

// NewVersion returns a new version with random values
func NewVersion(t *testing.T, a *articles.Article, v *articles.Version) *articles.Version {
	if v == nil {
		v = &articles.Version{}
	}

	if v.Title == "" {
		v.Title = uniuri.New()
	}

	if v.Subtitle == "" {
		v.Subtitle = uniuri.New()
	}

	if v.Description == "" {
		v.Description = uniuri.New()
	}

	if v.Content == "" {
		v.Content = uniuri.New()
	}

	v.ArticleID = a.ID

	if err := v.Create(); err != nil {
		t.Fatalf("failed to save article: %s", err)
	}

	lifecycle.SaveModels(t, v)
	return v
}
