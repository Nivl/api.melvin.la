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
	return nil
}

// SetAttributes sets the attributes of the file
func (s *Noop) SetAttributes(filepath string, attrs *UpdatableFileAttributes) (*FileAttributes, error) {
	return NewFileAttributesFromUpdatable(attrs), nil
}

// Attributes returns the attributes of the file
// Always returns an empty struct as no attributes are kept for this FS
func (s *Noop) Attributes(filepath string) (*FileAttributes, error) {
	return &FileAttributes{}, nil
}

// URL returns the URL of the file
func (s *Noop) URL(filepath string) (string, error) {
	return "http://localhost/noop/" + filepath, nil
}

// Exists check if a file exists
func (s *Noop) Exists(filepath string) (bool, error) {
	return true, nil
}

// Delete removes a file, ignores files that do not exist
func (s *Noop) Delete(filepath string) error {
	return nil
}

// WriteIfNotExist copies the provided io.Reader to dest if the file does
// not already exist
// Returns:
//   - A boolean specifying if the file got uploaded (true) or if already
//     existed (false).
//   - A URL to the uploaded file
//   - An error if something went wrong
func (s *Noop) WriteIfNotExist(src io.Reader, destPath string) (new bool, url string, err error) {
	url, _ = s.URL(destPath)
	return true, url, nil
}
