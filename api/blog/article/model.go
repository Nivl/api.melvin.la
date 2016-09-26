package article

import (
	"time"

	"github.com/Nivl/api.melvin.la/api/app"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func Query() *mgo.Collection {
	return app.GetContext().DB.C("article")
}

// Article is a structure representing an article that can be saved in the database
type Article struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Title       string        `bson:"title" json:"title"`
	Content     string        `bson:"content" json:"content"`
	Slug        string        `bson:"slug" json:"slug"`
	Subtitle    string        `bson:"subtitle" json:"subtitle"`
	Description string        `bson:"description" json:"description"`
	CreatedAt   time.Time     `bson:"createdAt" json:"created_at"`
}
