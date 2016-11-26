package articlestest

import (
	"testing"

	"github.com/dchest/uniuri"
	"github.com/gosimple/slug"
	"github.com/melvin-laplanche/ml-api/src/auth"
	"github.com/melvin-laplanche/ml-api/src/auth/authtest"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles"
	"github.com/melvin-laplanche/ml-api/src/db"
	"github.com/melvin-laplanche/ml-api/src/testhelpers"
)

// NewArticle returns a published article with random values
func NewArticle(t *testing.T, a *articles.Article) (*articles.Article, *auth.User, *auth.Session) {
	if a == nil {
		a = &articles.Article{
			PublishedAt: db.Now(),
		}
	}

	if a.Content == nil {
		a.Content = &articles.Content{}
	}

	if a.Content.Title == "" {
		a.Content.Title = uniuri.New()
	}

	if a.Slug == "" {
		a.Slug = slug.Make(a.Content.Title)
	}

	var user *auth.User
	var session *auth.Session

	if a.UserID == "" {
		user, session = authtest.NewAuth(t)
		a.User = user
		a.UserID = user.ID
	}

	if err := a.Create(); err != nil {
		t.Fatalf("failed to save article: %s", err)
	}

	a.Content.ArticleID = a.ID
	a.Content.IsCurrent = true
	if err := a.Content.Create(); err != nil {
		t.Fatalf("failed to save article content: %s", err)
	}

	if a.Draft != nil {
		a.Draft.ArticleID = a.ID
		a.Draft.IsDraft = true
		if err := a.Draft.Create(); err != nil {
			t.Fatalf("failed to save article draft: %s", err)
		}
	}

	testhelpers.SaveModels(t, a)

	return a, user, session
}
