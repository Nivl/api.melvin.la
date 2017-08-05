package filetype

import (
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
)

// ImageDecoder is a type that represents an image decoder
type ImageDecoder func(r io.Reader) (image.Image, error)

// IsImage checks the specified reader contains an image
func IsImage(r io.ReadSeeker) (bool, string, error) {
	mimeType, err := MimeType(r)
	if err != nil {
		return false, "", err
	}

	validType := map[string]FileValidator{
		"image/jpeg": IsJPG,
		"image/png":  IsPNG,
		"image/gif":  IsGIF,
	}
	validator, found := validType[mimeType]
	if !found {
		return false, "", nil
	}
	isValid, err := validator(r)
	return isValid, mimeType, err
}

// IsPNG validates a PNG file
func IsPNG(r io.ReadSeeker) (bool, error) {
	return validateImage(r, png.Decode)
}

// IsJPG validates a JPG file
func IsJPG(r io.ReadSeeker) (bool, error) {
	return validateImage(r, jpeg.Decode)
}

// IsGIF validates a GIF file
func IsGIF(r io.ReadSeeker) (bool, error) {
	return validateImage(r, gif.Decode)
}

// validateImage check if an images has a valid format
// Update when this gets done: https://github.com/golang/go/issues/18098
func validateImage(r io.ReadSeeker, decode ImageDecoder) (bool, error) {
	initialPos, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}

	// Parse the whole file
	_, err = decode(r)
	success := (err == nil)

	// revert the pointer back to its original position
	_, err = r.Seek(initialPos, io.SeekStart)
	if err != nil {
		return false, err
	}
	return success, nil
}
