package mockfilestorage

import (
	"os"
	"path"

	"github.com/Nivl/go-rest-tools/storage/filestorage"
	"github.com/stretchr/testify/mock"
)

var (
	// StringType represent a string argument
	StringType = mock.AnythingOfType("string")

	// UpdatableAttrType represents a *filestorage.UpdatableFileAttributes argument
	UpdatableAttrType = mock.AnythingOfType("*filestorage.UpdatableFileAttributes")
)

// ExpectWriteIfNotExist is an helper that expects a WriteIfNotExist to
// succeed with the provided params
func (s *FileStorage) ExpectWriteIfNotExist(isNew bool, url string) *mock.Call {
	return s.On("WriteIfNotExist", mock.Anything, StringType).Return(isNew, url, nil)
}

// ExpectRead is an helper that expects a Read
func (s *FileStorage) ExpectRead(cwd, filename string) *mock.Call {
	filePath := path.Join(cwd, "fixtures", filename)
	return s.On("Read", StringType).Return(os.Open(filePath))
}

// ExpectExists is an helper that expects Exists() to return true
func (s *FileStorage) ExpectExists() *mock.Call {
	return s.On("Exists", StringType).Return(true, nil)
}

// ExpectNotExists is an helper that expects Exists() to return false
func (s *FileStorage) ExpectNotExists() *mock.Call {
	return s.On("Exists", StringType).Return(false, nil)
}

// ExpectURL is an helper that expects URL() to return given param
func (s *FileStorage) ExpectURL(url string) *mock.Call {
	return s.On("URL", StringType).Return(url, nil)
}

// ExpectSetAttributes is an helper that expects SetAttributes() to work,
// and returns an empty content
func (s *FileStorage) ExpectSetAttributes() *mock.Call {
	attrs := &filestorage.FileAttributes{}
	return s.On("SetAttributes", StringType, UpdatableAttrType).Return(attrs, nil)
}

// ExpectSetAttributesRet is an helper that expects SetAttributes() to
// return the provided object
func (s *FileStorage) ExpectSetAttributesRet(attrs *filestorage.FileAttributes) *mock.Call {
	return s.On("SetAttributes", StringType, UpdatableAttrType).Return(attrs, nil)
}

// ExpectAttributes is an helper that expects Attributes() to
// return the provided object
func (s *FileStorage) ExpectAttributes(attrs *filestorage.FileAttributes) *mock.Call {
	return s.On("Attributes", StringType).Return(attrs, nil)
}

// ExpectDelete is an helper that expects Delete() to succeed
func (s *FileStorage) ExpectDelete() *mock.Call {
	return s.On("Delete", StringType).Return(nil)
}
