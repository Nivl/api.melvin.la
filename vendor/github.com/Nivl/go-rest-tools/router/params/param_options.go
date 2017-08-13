package params

import (
	"io"
	"mime/multipart"
	"reflect"
	"strconv"
	"strings"

	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-rest-tools/types/filetype"
	"github.com/Nivl/go-rest-tools/types/slices"
	"github.com/Nivl/go-rest-tools/types/strngs"
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

	// ErrMsgMaxLen represents the error message corresponding to
	// a field that exceed the maximum number of char
	ErrMsgMaxLen = "too many chars"

	// ErrMsgEnum represents the error message corresponding to
	// a field that doesn't contain a value set in an enum
	ErrMsgEnum = "not a valid value"
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

	// MaxLen represents the maximum length a param can have (under its string
	// form). Any invalid values (including 0) will be ignored
	// maxlen:"255"
	MaxLen int

	// AuthorizedValues represents the list of authorized value for this param
	// enum:"and,or"
	AuthorizedValues []string
}

// Validate checks the given value passes the options set
func (opts *ParamOptions) Validate(value string, wasProvided bool) error {
	if opts.MaxLen > 0 && len(value) > opts.MaxLen {
		return apierror.NewBadRequest(opts.Name, ErrMsgMaxLen)
	}

	if value == "" && opts.Required {
		return apierror.NewBadRequest(opts.Name, ErrMsgMissingParameter)
	}

	if value == "" && opts.NoEmpty && wasProvided {
		return apierror.NewBadRequest(opts.Name, ErrMsgEmptyParameter)
	}

	if value != "" {
		if opts.ValidateUUID && !strngs.IsValidUUID(value) {
			return apierror.NewBadRequest(opts.Name, ErrMsgInvalidUUID)
		}

		if opts.ValidateURL && !strngs.IsValidURL(value) {
			return apierror.NewBadRequest(opts.Name, ErrMsgInvalidURL)
		}

		if opts.ValidateEmail && !strngs.IsValidEmail(value) {
			return apierror.NewBadRequest(opts.Name, ErrMsgInvalidEmail)
		}

		if len(opts.AuthorizedValues) > 0 {
			found, err := slices.InSlice(opts.AuthorizedValues, value)
			if err != nil {
				return err
			}
			if !found {
				return apierror.NewBadRequest(opts.Name, ErrMsgEnum)
			}
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
		return "", apierror.NewBadRequest(opts.Name, errorMsg)
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

	// We use the maxlen tag to get the max length of a the value
	maxlen := tags.Get("maxlen")
	if len(maxlen) > 0 {
		// we silently fail on errors
		output.MaxLen, _ = strconv.Atoi(maxlen)
	}

	// We use the enu, tag to get all the authorized value a param can have
	enum := tags.Get("enum")
	if len(enum) > 0 {
		// we silently fail on errors
		output.AuthorizedValues = strings.Split(enum, ",")
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
