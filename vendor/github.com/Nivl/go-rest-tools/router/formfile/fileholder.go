package formfile

import "mime/multipart"

// FileHolder represents an interface to fetch files from a key
type FileHolder interface {
	FormFile(key string) (multipart.File, *multipart.FileHeader, error)
}
