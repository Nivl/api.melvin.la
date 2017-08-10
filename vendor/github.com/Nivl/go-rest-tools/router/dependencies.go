package router

import (
	"context"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/notifiers/mailer"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/Nivl/go-rest-tools/storage/filestorage"
)

// Dependencies represents all the dependencies of the API
type Dependencies struct {
	DB      db.DB
	Mailer  mailer.Mailer
	Storage filestorage.FileStorage
}

// NewDefaultDependencies returns the default dependencies using a TODO
// context
func NewDefaultDependencies() (*Dependencies, error) {
	return NewDefaultDependenciesWithContext(context.TODO())
}

// NewDefaultDependenciesWithContext returns the defaults dependencies using
// the provided context
func NewDefaultDependenciesWithContext(ctx context.Context) (*Dependencies, error) {
	storage, err := dependencies.NewStorage(ctx)
	if err != nil {
		return nil, err
	}

	return &Dependencies{
		DB:      dependencies.DB,
		Mailer:  dependencies.NewMailer(),
		Storage: storage,
	}, nil
}

// NewNoFailersDependencies returns the default dependencies that cannot
// fails during creation and replaces the others by noops
func NewNoFailersDependencies() *Dependencies {
	return &Dependencies{
		DB:      dependencies.DB,
		Mailer:  dependencies.NewMailer(),
		Storage: &filestorage.Noop{},
	}
}
