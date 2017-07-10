package filestorage

import (
	"errors"
	"io"
)

// Noop is an implementation of the FileStorage interface that does nothing
type Noop struct {
}

// ID returns the unique identifier of the storage provider
func (s *Noop) ID() string {
	return "noop"
}

// SetBucket is used to set the bucket
// Always return nil
func (s *Noop) SetBucket(name string) error {
	return nil
}

// Read fetches a file and returns a reader
func (s *Noop) Read(filepath string) (io.ReadCloser, error) {
	return nil, errors.New("noop cannot read")
}

// Write copy the provided os.File to dest
func (s *Noop) Write(src io.Reader, destPath string) error {
	return errors.New("noop cannot write")
}

// SetAttributes sets the attributes of the file
func (s *Noop) SetAttributes(filepath string, attrs *UpdatableFileAttributes) (*FileAttributes, error) {
	return &FileAttributes{
		ContentType:        attrs.ContentType.(string),
		ContentDisposition: attrs.ContentDisposition.(string),
		ContentLanguage:    attrs.ContentLanguage.(string),
		ContentEncoding:    attrs.ContentEncoding.(string),
		CacheControl:       attrs.CacheControl.(string),
		Metadata:           attrs.Metadata,
	}, nil
}

// Attributes returns the attributes of the file
// Always returns an empty struct as no attributes are kept for this FS
func (s *Noop) Attributes(filepath string) (*FileAttributes, error) {
	return &FileAttributes{}, nil
}

// URL returns the URL of the file
func (s *Noop) URL(filepath string) (string, error) {
	return "", errors.New("noop cannot generate a URL")
}

// Delete removes a file, ignores files that do not exist
func (s *Noop) Delete(filepath string) error {
	return errors.New("noop cannot read")
}
