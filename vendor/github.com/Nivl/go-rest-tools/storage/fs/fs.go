package fs

import "os"

// Exists checks if a file or a directory exists
func Exists(filepath string) (bool, error) {
	_, err := os.Stat(filepath)
	if err == nil {
		return true, err
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FileExists checks if a file exists
func FileExists(filepath string) (bool, error) {
	fi, err := os.Stat(filepath)
	if err == nil {
		return !fi.IsDir(), err
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// DirExists checks if a directory exists
func DirExists(filepath string) (bool, error) {
	fi, err := os.Stat(filepath)
	if err == nil {
		return fi.IsDir(), err
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
