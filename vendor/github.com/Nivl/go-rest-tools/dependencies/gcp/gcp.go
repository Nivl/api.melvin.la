package gcp

import (
	"context"

	"github.com/Nivl/go-rest-tools/storage/filestorage"
)

// GCP represents a Google Cloud Platform service
type GCP interface {
	// Storage creates a new storage attached to the provided context
	Storage(ctx context.Context) (filestorage.FileStorage, error)
}
