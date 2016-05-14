package article

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Article is a structure representing an article that can be saved in the database
type Article struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Title       string        `bson:"title" json:"title" form:"title" binding:"required"`
	Content     string        `bson:"content" json:"content" form:"content" binding:"required"`
	Slug        string        `bson:"slug" json:"slug" form:"slug" binding:"required"`
	Subtitle    string        `bson:"subtitle" json:"subtitle" form:"subtitle" binding:"required"`
	Description string        `bson:"description" json:"description" form:"description" binding:"required"`
	CreatedAt   time.Time     `bson:"createdAt" json:"created_at"`
}

// ToCollection returns an []interface{} from []Article
func ToCollection(list []Article) []interface{} {
	slice := make([]interface{}, len(list))

	for i, d := range list {
		slice[i] = d
	}

	return slice
}
