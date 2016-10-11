package articles

import (
	"github.com/Nivl/api.melvin.la/api/app"
	"gopkg.in/mgo.v2"
)

// EnsureIndexes sets the indexes for the Articles document
func EnsureIndexes() {
	indexes := []mgo.Index{
		mgo.Index{Key: []string{"slug"}, Unique: true, DropDups: true, Background: true},
		mgo.Index{Key: []string{"-created_at"}, Background: true},
	}
	doc := app.GetContext().DB.C("articles")

	for _, index := range indexes {
		if err := doc.EnsureIndex(index); err != nil {
			panic(err)
		}
	}
}
