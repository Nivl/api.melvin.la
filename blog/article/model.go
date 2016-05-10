package article

import "gopkg.in/mgo.v2/bson"

// Article is a structure representing an article that can be saved in the database
type Article struct {
	ID      bson.ObjectId `bson:"_id" json:"id"`
	Title   string        `bson:"title" json:"title"`
	Content string        `bson:"content" json:"content"`
}

// ToCollection returns an []interface{} from []Article
func ToCollection(list []Article) []interface{} {
	slice := make([]interface{}, len(list))

	for i, d := range list {
		slice[i] = d
	}

	return slice
}
