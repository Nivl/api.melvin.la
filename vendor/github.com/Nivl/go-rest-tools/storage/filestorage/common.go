package filestorage

import "io"

// writeIfNotExist is an helper to implements FileStorage.WriteIfNotExist
// Since this function is the same for all provider, let's not rewrite it
func writeIfNotExist(s FileStorage, src io.Reader, destPath string) (new bool, url string, err error) {
	exists, err := s.Exists(destPath)
	if err != nil {
		return false, "", err
	}

	if !exists {
		if err = s.Write(src, destPath); err != nil {
			return false, "", err
		}
	}

	url, err = s.URL(destPath)
	if err != nil {
		return false, "", err
	}

	return !exists, url, nil
}
