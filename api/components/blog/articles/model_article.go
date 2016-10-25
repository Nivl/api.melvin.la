package articles

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/app"
	"github.com/dchest/uniuri"
	"github.com/gosimple/slug"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Query() *mgo.Collection {
	return app.GetContext().DB.C("articles")
}

var defaultSearch = bson.M{
	"is_deleted": false,
}

// Article is a structure representing an article that can be saved in the database
type Article struct {
	ID          bson.ObjectId `bson:"_id"`
	Title       string        `bson:"title"`
	Content     string        `bson:"content"`
	Slug        string        `bson:"slug"`
	Subtitle    string        `bson:"subtitle"`
	Description string        `bson:"description"`
	CreatedAt   time.Time     `bson:"created_at"`
	IsDeleted   bool          `bson:"is_deleted"`
	IsPublished bool          `bson:"is_published"`
}

func (a *Article) FullyDelete() error {
	if a == nil {
		return errors.New("article not instanced")
	}

	if a.ID == "" {
		return errors.New("article has not been saved")
	}

	return Query().RemoveId(a.ID)
}

func (a *Article) Save() error {
	if a == nil {
		return errors.New("article not instanced")
	}

	if a.ID == "" {
		return a.Create()
	}

	return a.Update()
}

func (a *Article) Create() error {
	if a == nil {
		return apierror.NewServerError("article not instanced")
	}

	if a.Slug == "" {
		a.Slug = slug.Make(a.Title)
	}

	if bson.IsObjectIdHex(a.Slug) {
		return apierror.NewBadRequest("slug cannot be a ObjectId")
	}

	a.CreatedAt = time.Now()

	// To prevent duplicates on the slug, we'll retry the insert() up to 10 times
	originalSlug := a.Slug
	var err error
	for i := 0; i < 10; i++ {
		a.ID = bson.NewObjectId()
		err = Query().Insert(a)

		if err != nil {
			// In case of duplicate we'll add "-X" at the end of the slug, where X is
			// a number
			a.Slug = fmt.Sprintf("%s-%d", originalSlug, i)

			if mgo.IsDup(err) == false {
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

func (a *Article) Update() error {
	return nil
}

func NewTestArticle(t *testing.T, a *Article) *Article {
	if a == nil {
		a = &Article{
			IsDeleted:   false,
			IsPublished: true,
		}
	}

	if a.Title == "" {
		a.Title = uniuri.New()
	}

	if err := a.Save(); err != nil {
		t.Fatalf("failed to save article: %s", err)
	}
	return a
}
