package router

import (
	"github.com/Nivl/go-rest-tools/notifiers/mailer"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/storage/filestorage"
)

// Dependencies represents all the dependencies of the API
type Dependencies struct {
	DB      db.Connection
	Mailer  mailer.Mailer
	Storage filestorage.FileStorage
}
