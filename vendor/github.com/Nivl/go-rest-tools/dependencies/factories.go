package dependencies

import (
	"context"

	"github.com/Nivl/go-rest-tools/logger"
	"github.com/Nivl/go-rest-tools/notifiers/mailer"
	"github.com/Nivl/go-rest-tools/storage/filestorage"
)

// NewLogger returns a Logger
func NewLogger() logger.Logger {
	if Logentries != nil {
		return logger.NewLogEntries(Logentries)
	}
	return logger.NewBasicLogger()
}

// NewMailer returns a Mailer
func NewMailer() mailer.Mailer {
	if Sendgrid != nil {
		return mailer.NewSendgrid(Sendgrid.APIKey, Sendgrid.From, Sendgrid.To, Sendgrid.StacktraceUUID)
	}
	return &mailer.Noop{}
}

// NewStorage returns a file storage provide
func NewStorage(ctx context.Context) (filestorage.FileStorage, error) {
	var storage filestorage.FileStorage
	var err error
	bucket := "ml-api"

	if GoogleCloud != nil {
		storage, err = filestorage.NewGCStorage(ctx, GoogleCloud.APIKey)
		bucket = GoogleCloud.Bucket
	} else if Cloudinary != nil {
		storage = filestorage.NewCloudinary(Cloudinary.APIKey, Cloudinary.Secret)
		bucket = Cloudinary.Bucket
	}

	if err != nil {
		return nil, err
	}
	if err = storage.SetBucket(bucket); err != nil {
		return nil, err
	}
	return storage, nil
}
