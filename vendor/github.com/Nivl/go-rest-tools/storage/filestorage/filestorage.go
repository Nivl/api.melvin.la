package filestorage

import (
	"io"
)

// FileStorage is a interface to store and retrieve files
type FileStorage interface {
	// ID returns the unique identifier of the storage provider
	ID() string

	// SetBucket sets the bucket that will contains the files
	SetBucket(bucket string) error

	// Read fetches a file a returns a reader
	// Returns os.ErrNotExist if the file does not exists
	Read(filepath string) (io.ReadCloser, error)

	// Write copy the provided os.File to dest
	Write(src io.Reader, destPath string) error

	// Delete removes a file, ignores files that do not exist
	Delete(filepath string) error

	// URL returns the URL of the file
	URL(filepath string) (string, error)

	// SetAttributes sets the attributes of the file
	SetAttributes(filepath string, attrs *UpdatableFileAttributes) (*FileAttributes, error)

	// Attributes returns the attributes of the file
	Attributes(filepath string) (*FileAttributes, error)
}

// FileAttributes represents the attributes a file can have
type FileAttributes struct {
	ContentType        string
	ContentLanguage    string
	ContentEncoding    string
	ContentDisposition string
	CacheControl       string
	Metadata           map[string]string
}

// UpdatableFileAttributes represents the updatable attributes a file can have
type UpdatableFileAttributes struct {
	ContentType        interface{}
	ContentLanguage    interface{}
	ContentEncoding    interface{}
	ContentDisposition interface{}
	CacheControl       interface{}
	Metadata           map[string]string // set to map[string]string{} to delete
}
