package dependencies

import (
	"context"

	"github.com/Nivl/go-rest-tools/dependencies/gcp"
	"github.com/Nivl/go-rest-tools/logger"
	"github.com/Nivl/go-rest-tools/notifiers/mailer"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/storage/db/sqlx"
	"github.com/Nivl/go-rest-tools/storage/filestorage"
	"github.com/bsphere/le_go"
)

var _ Dependencies = (*APIDependencies)(nil)

// APIDependencies is a structure implementing the Dependencies interface
type APIDependencies struct {
	postgres   *sqlx.Connection
	logentries *le_go.Logger
	sendgrid   mailer.Mailer
	gcp        gcp.GCP
	cloudinary filestorage.FileStorage

	logger logger.Logger
	mailer mailer.Mailer
}

// SetDB creates a connection to a SQL database
func (deps *APIDependencies) SetDB(uri string) error {
	var err error
	deps.postgres, err = sqlx.New(uri)
	return err
}

// DB returns the current SQL connection
func (deps *APIDependencies) DB() db.Connection {
	return deps.postgres
}

// SetLogentries creates a connection to logentries
func (deps *APIDependencies) SetLogentries(token string) error {
	var err error
	deps.logentries, err = le_go.Connect(token)
	return err
}

// Logger returns the default logger following this order:
// Logentries
// BasicLogger
func (deps *APIDependencies) Logger() logger.Logger {
	if deps.logger == nil {
		if deps.logentries != nil {
			deps.logger = logger.NewLogEntries(deps.logentries)
		} else {
			deps.logger = logger.NewBasicLogger()
		}
	}
	return deps.logger
}

// SetSendgrid creates a mailer using sendgrid
func (deps *APIDependencies) SetSendgrid(apiKey, from, to, stacktraceUUID string) error {
	deps.sendgrid = mailer.NewSendgrid(apiKey, from, to, stacktraceUUID)
	return nil
}

// Mailer returns the default mailer following this order:
// Sendgrid
// Noop
func (deps *APIDependencies) Mailer() mailer.Mailer {
	if deps.mailer == nil {
		if deps.sendgrid != nil {
			deps.mailer = deps.sendgrid
		} else {
			deps.mailer = &mailer.Noop{}
		}
	}
	return deps.mailer
}

// SetGCP sets up Google Cloud Platform
func (deps *APIDependencies) SetGCP(apiKey, projectName, bucket string) error {
	deps.gcp = gcp.New(apiKey, projectName, bucket)
	return nil
}

// SetCloudinary setups Cloudinary as Storage provider
func (deps *APIDependencies) SetCloudinary(apiKey, secret, bucket string) error {
	deps.cloudinary = filestorage.NewCloudinary(apiKey, secret)
	deps.cloudinary.SetBucket(bucket)
	return nil
}

// FileStorage returns the default filestorage following this order
// GCP
// Cloudinary
// Filesystem
func (deps *APIDependencies) FileStorage(ctx context.Context) (filestorage.FileStorage, error) {
	if deps.gcp != nil {
		return deps.gcp.Storage(ctx)
	}
	if deps.cloudinary != nil {
		return deps.cloudinary, nil
	}
	return filestorage.NewFSStorage()
}
