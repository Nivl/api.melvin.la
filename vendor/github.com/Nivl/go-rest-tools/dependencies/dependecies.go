package dependencies

import (
	"context"

	"github.com/Nivl/go-rest-tools/logger"
	"github.com/Nivl/go-rest-tools/notifiers/mailer"
	"github.com/Nivl/go-rest-tools/notifiers/reporter"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/storage/filestorage"
)

// Dependencies represents the dependency of the api
type Dependencies interface {
	// SetDB creates a connection to a SQL database
	SetDB(uri string) error

	// DB returns the current SQL connection
	DB() db.Connection

	// SetLogentries creates a connection to logentries
	SetLogentries(token string) error

	// Logger returns the default logger following this order:
	// Logentries
	// BasicLogger
	Logger() logger.Logger

	// SetSendgrid creates a mailer using sendgrid
	SetSendgrid(apiKey, from, to, stacktraceUUID string) error

	// Mailer returns the default mailer following this order:
	// Sendgrid
	// Noop
	Mailer() mailer.Mailer

	// SetGCP sets up Google Cloud Platform
	SetGCP(apiKey, projectName, bucket string) error

	// SetCloudinary setups Cloudinary as Storage provider
	SetCloudinary(apiKey, secret, bucket string) error

	// FileStorage returns the default filestorage following this order
	// GCP
	// Cloudinary
	// Filesystem
	FileStorage(ctx context.Context) (filestorage.FileStorage, error)

	// SetSentry creates a reporter using Sentry
	SetSentry(con string) error

	// EnableEmailReporting sets the current mailer as reporter
	EnableEmailReporting(con string) error

	// Reporter returns the default reporter creator following this order:
	// Sentry
	// Email
	// Noop
	Reporter() reporter.Creator
}
