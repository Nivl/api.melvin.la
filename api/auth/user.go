package auth

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/Nivl/api.melvin.la/api/apierror"
	"github.com/Nivl/api.melvin.la/api/app"
	"github.com/dchest/uniuri"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func QueryUsers() *mgo.Collection {
	return app.GetContext().DB.C("users")
}

// User is a structure representing a user that can be saved in the database
type User struct {
	ID        bson.ObjectId `bson:"_id"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
	IsDeleted bool          `bson:"is_deleted"`

	Name         string    `bson:"name"`
	Email        string    `bson:"email"`
	Password     string    `bson:"password"`
	SignUpIp     string    `bson:"sign_up_ip"`
	LastActionIp string    `bson:"last_action_ip"`
	LastActionAt time.Time `bson:"last_action_at"`
}

func GetUser(userID bson.ObjectId) (*User, error) {
	if !userID.Valid() {
		return nil, apierror.NewServerError("not a valid user id")
	}

	toFind := bson.M{
		"_id":        userID,
		"is_deleted": false,
	}

	var user User
	err := QueryUsers().Find(toFind).One(&user)
	return &user, err
}

func CryptPassword(raw string) (string, error) {
	password, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func (u *User) Save() error {
	if u == nil {
		return apierror.NewServerError("user is not instanced")
	}

	return u.Create()
}

func (u *User) FullyDelete() error {
	if u == nil {
		return errors.New("user not instanced")
	}

	if u.ID == "" {
		return errors.New("user has not been saved")
	}

	return QueryUsers().RemoveId(u.ID)
}

func (u *User) Create() error {
	if u == nil {
		return apierror.NewServerError("user is not instanced")
	}

	if u.ID != "" {
		return apierror.NewServerError("cannot persist a user that has an ID")
	}

	u.ID = bson.NewObjectId()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	err := QueryUsers().Insert(u)
	if err != nil && mgo.IsDup(err) {
		return apierror.NewConflict("email address already in use")
	}

	return err
}

func NewTestUser(t *testing.T, u *User) *User {
	if u == nil {
		u = &User{
			IsDeleted: false,
		}
	}

	if u.Email == "" {
		u.Email = fmt.Sprintf("fake+%s@melvin.la", uniuri.New())
	}

	if u.Name == "" {
		u.Email = "Fake Account"
	}

	if u.Password == "" {
		var err error
		u.Password, err = CryptPassword("fake")
		if err != nil {
			t.Fatalf("failed to create password: %s", err)
		}
	}

	if err := u.Create(); err != nil {
		t.Fatalf("failed to create user: %s", err)
	}
	return u
}
