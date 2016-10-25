package auth

import (
	"testing"

	"github.com/Nivl/api.melvin.la/api/app"
	mgo "gopkg.in/mgo.v2"
)

// EnsureIndexes sets the indexes for the Users and Sessions document
func EnsureIndexes() {
	indexes := []mgo.Index{
		mgo.Index{Key: []string{"email"}, Unique: true, Background: true},
		mgo.Index{Key: []string{"-created_at"}, Background: true},
	}
	doc := app.GetContext().DB.C("users")

	for _, index := range indexes {
		if err := doc.EnsureIndex(index); err != nil {
			panic(err)
		}
	}
}

func NewTestAuth(t *testing.T) (*User, *Session) {
	user := NewTestUser(t, nil)
	session := &Session{
		UserID: user.ID,
	}

	if err := session.Create(); err != nil {
		t.Fatal(err)
	}

	return user, session
}
