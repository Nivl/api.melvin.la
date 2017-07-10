package filestorage

import (
	"context"
	"io"

	"google.golang.org/api/option"

	"fmt"

	"cloud.google.com/go/storage"
)

// toGCStorage converts a *UpdatableFileAttributes into a *storage.ObjectAttrsToUpdate
func (attrs *UpdatableFileAttributes) toGCStorage() *storage.ObjectAttrsToUpdate {
	return &storage.ObjectAttrsToUpdate{
		ContentType:        attrs.ContentType,
		ContentDisposition: attrs.ContentDisposition,
		ContentLanguage:    attrs.ContentLanguage,
		ContentEncoding:    attrs.ContentEncoding,
		CacheControl:       attrs.CacheControl,
		Metadata:           attrs.Metadata,
	}
}

// NewAttributesFromGCStorage converts a *storage.ObjectAttrs to *FileAttributes
func NewAttributesFromGCStorage(attrs *storage.ObjectAttrs) *FileAttributes {
	return &FileAttributes{
		ContentType:        attrs.ContentType,
		ContentDisposition: attrs.ContentDisposition,
		ContentLanguage:    attrs.ContentLanguage,
		ContentEncoding:    attrs.ContentEncoding,
		CacheControl:       attrs.CacheControl,
		Metadata:           attrs.Metadata,
	}
}

// NewGCStorage returns a new instance of a Google Cloud Storage
func NewGCStorage(ctx context.Context, apiKey string) (*GCStorage, error) {
	client, err := storage.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return &GCStorage{
		ctx:    ctx,
		client: client,
	}, nil
}

// GCStorage is an implementation of the FileStorage interface for Google Cloud
type GCStorage struct {
	ctx        context.Context
	client     *storage.Client
	bucket     *storage.BucketHandle
	bucketName string
}

// ID returns the unique identifier of the storage provider
func (s *GCStorage) ID() string {
	return "google_cloud_storage"
}

// SetBucket is used to set the bucket
// Always return nil
func (s *GCStorage) SetBucket(name string) error {
	s.bucketName = name
	s.bucket = s.client.Bucket(name)
	return nil
}

// Read fetches a file a returns a reader
func (s *GCStorage) Read(filepath string) (io.ReadCloser, error) {
	return s.bucket.Object(filepath).NewReader(s.ctx)
}

// Write copy the provided os.File to dest
func (s *GCStorage) Write(src io.Reader, destPath string) error {
	obj := s.bucket.Object(destPath)
	dest := obj.NewWriter(s.ctx)

	// Copy the file
	_, err := io.Copy(dest, src)
	if err != nil {
		dest.CloseWithError(err)
		return err
	}

	// Send the changes to GCP
	if err := dest.Close(); err != nil {
		return err
	}

	return nil
}

// SetAttributes sets the attributes of the file
func (s *GCStorage) SetAttributes(filepath string, attrs *UpdatableFileAttributes) (*FileAttributes, error) {
	gcsAttrs, err := s.bucket.Object(filepath).Update(s.ctx, *attrs.toGCStorage())
	if err != nil {
		return nil, err
	}

	return NewAttributesFromGCStorage(gcsAttrs), nil
}

// Attributes returns the attributes of the file
func (s *GCStorage) Attributes(filepath string) (*FileAttributes, error) {
	gcsAttrs, err := s.bucket.Object(filepath).Attrs(s.ctx)
	if err != nil {
		return nil, err
	}
	return NewAttributesFromGCStorage(gcsAttrs), nil
}

// URL returns the URL of the file
func (s *GCStorage) URL(filepath string) (string, error) {
	return fmt.Sprintf("https://%s.storage.googleapis.com/%s", s.bucketName, filepath), nil
}

// Delete removes a file, ignores files that do not exist
func (s *GCStorage) Delete(filepath string) error {
	return s.bucket.Object(filepath).Delete(s.ctx)
}
