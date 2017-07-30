package params

import (
	"io"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/primitives/filetype"
	"github.com/Nivl/go-rest-tools/primitives/strngs"
)

const (
	// ErrMsgMissingParameter represents the error message corresponding to
	// a missing param
	ErrMsgMissingParameter = "parameter missing"

	// ErrMsgEmptyParameter represents the error message corresponding to
	// a missing param
	ErrMsgEmptyParameter = "parameter can be omitted but not empty"

	// ErrMsgInvalidUUID represents the error message corresponding to
	// an invalid UUID
	ErrMsgInvalidUUID = "not a valid uuid"

	// ErrMsgInvalidURL represents the error message corresponding to
	// an invalid URL
	ErrMsgInvalidURL = "not a valid url"

	// ErrMsgInvalidEmail represents the error message corresponding to
	// an invalid Email address
	ErrMsgInvalidEmail = "not a valid email"

	// ErrMsgInvalidImage represents the error message corresponding to
	// an invalid image
	ErrMsgInvalidImage = "not a valid image"
)

// ParamOptions represent all the options for a field
type ParamOptions struct {
	// Ignore means the field should not been parsed
	// json:"-"
	Ignore bool

	// Name contains the name of the field in the payload
	// json:"my_field"
	Name string

	// Required means the request should fail with a Bad Request if the field is missing.
	// params:"required"
	Required bool

	// Trim means the field needs to be trimmed before being retrieved and checked
	// params:"trim"
	Trim bool

	// ValidateUUID means the field should contain a valid UUIDv4
	// params:"uuid"
	ValidateUUID bool

	// ValidateEmail means the field should contain a valid email
	// params:"email"
	ValidateEmail bool

	// ValidateURL means the field should contain a valid url
	// params:"url"
	ValidateURL bool

	// ValidateImage means the field should contain a valid image
	// params:"image"
	ValidateImage bool

	// NoEmpty means the field should not contain an empty value
	// The difference between `required` and `noempty` is that `require` does`
	// not accept nil pointer. `noempty` accepts nil pointer, but if a value is
	// provided it cannot be an empty string.
	// params:"noempty"
	NoEmpty bool
}

// Validate checks the given value passes the options set
func (opts *ParamOptions) Validate(value string, wasProvided bool) error {
	if value == "" && opts.Required {
		return httperr.NewBadRequest(opts.Name, ErrMsgMissingParameter)
	}

	if value == "" && opts.NoEmpty && wasProvided {
		return httperr.NewBadRequest(opts.Name, ErrMsgEmptyParameter)
	}

	if value != "" {
		if opts.ValidateUUID && !strngs.IsValidUUID(value) {
			return httperr.NewBadRequest(opts.Name, ErrMsgInvalidUUID)
		}

		if opts.ValidateURL && !strngs.IsValidURL(value) {
			return httperr.NewBadRequest(opts.Name, ErrMsgInvalidURL)
		}

		if opts.ValidateEmail && !strngs.IsValidEmail(value) {
			return httperr.NewBadRequest(opts.Name, ErrMsgInvalidEmail)
		}
	}

	return nil
}

// ValidateFileContent checks the given file passes the options set
func (opts *ParamOptions) ValidateFileContent(file multipart.File) (string, error) {
	// Just by security, but it shouldn't be necessary
	defer file.Seek(0, io.SeekStart)

	var valid bool
	var mime string
	var err error
	var errorMsg string

	if opts.ValidateImage {
		valid, mime, err = filetype.IsImage(file)
		errorMsg = ErrMsgInvalidImage
	} else {
		// We still get the mimetype
		valid = true
		mimeType, err := filetype.MimeType(file)
		if err != nil {
			return "", err
		}
		return mimeType, nil
	}

	if err != nil {
		return "", err
	}

	if !valid {
		return "", httperr.NewBadRequest(opts.Name, errorMsg)
	}

	// check "valid", and return an error if its not
	return mime, nil
}

// ApplyTransformations applies all the wanted transformations to the given value
func (opts *ParamOptions) ApplyTransformations(value string) string {
	if opts.Trim {
		value = strings.TrimSpace(value)
	}
	return value
}

// NewParamOptions returns a ParamOptions from a StructTag
func NewParamOptions(tags *reflect.StructTag) *ParamOptions {
	output := &ParamOptions{}

	// We use the json tag to get the field name
	jsonOpts := strings.Split(tags.Get("json"), ",")
	if len(jsonOpts) > 0 {
		if jsonOpts[0] == "-" {
			return &ParamOptions{Ignore: true}
		}

		output.Name = jsonOpts[0]
	}

	// We parse the params
	opts := strings.Split(tags.Get("params"), ",")
	nbOptions := len(opts)
	for i := 0; i < nbOptions; i++ {
		switch opts[i] {
		case "required":
			output.Required = true
		case "noempty":
			output.NoEmpty = true
		case "trim":
			output.Trim = true
		case "uuid":
			output.ValidateUUID = true
		case "email":
			output.ValidateEmail = true
		case "url":
			output.ValidateURL = true
		case "image":
			output.ValidateImage = true
		}
	}

	return output
}
