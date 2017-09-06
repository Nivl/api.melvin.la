package gcp

import (
	"context"

	"github.com/Nivl/go-rest-tools/storage/filestorage"
)

var _ GCP = (*Service)(nil)

// New creates a new Service instance
func New(apiKey, projectName, bucket string) GCP {
	return &Service{
		apiKey:      apiKey,
		projectName: projectName,
		bucket:      bucket,
	}
}

// Service implements the GCP interface
type Service struct {
	apiKey      string
	projectName string
	bucket      string

	storage filestorage.FileStorage
}

// Storage returns a new storage instance
func (service *Service) Storage(ctx context.Context) (filestorage.FileStorage, error) {
	storage, err := filestorage.NewGCStorage(ctx, service.apiKey)
	if err != nil {
		return nil, err
	}
	storage.SetBucket(service.bucket)
	return storage, nil
}
