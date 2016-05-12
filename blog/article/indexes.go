package article

import (
	"github.com/Nivl/api.melvin.la/app"
	"gopkg.in/mgo.v2"
)

// EnsureIndexes sets the indexes for the Article document
func EnsureIndexes() {
	indexes := []mgo.Index{
		mgo.Index{Key: []string{"slug"}, Unique: true, DropDups: true, Background: true},
		mgo.Index{Key: []string{"-createdAt"}, Background: true},
	}
	doc := app.GetContext().DB.C("article")

	for _, index := range indexes {
		if err := doc.EnsureIndex(index); err != nil {
			panic(err)
		}
	}
}
