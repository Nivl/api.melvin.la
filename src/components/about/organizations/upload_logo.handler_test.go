package organizations_test

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/formfile/mockformfile"
	"github.com/Nivl/go-rest-tools/router/formfile/testformfile"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/router/params"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/storage/filestorage/mockfilestorage"
	"github.com/melvin-laplanche/ml-api/src/components/about/organizations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadLogoAccess(t *testing.T) {
	testCases := []testguard.AccessTestCase{
		{
			Description: "Should fail for anonymous users",
			User:        nil,
			ErrCode:     http.StatusUnauthorized,
		},
		{
			Description: "Should fail for logged users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
			ErrCode:     http.StatusForbidden,
		},
		{
			Description: "Should work for admin users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0", IsAdmin: true},
			ErrCode:     0,
		},
	}

	g := organizations.Endpoints[organizations.EndpointUploadLogo].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestUploadLogoInvalidParams(t *testing.T) {

	// create the multipart data
	cwd, _ := os.Getwd()
	licenseHeader, licenseFile := testformfile.NewMultipartData(t, cwd, "LICENSE")
	defer licenseFile.Close()

	imageHeader, imageFile := testformfile.NewMultipartData(t, cwd, "black_pixel.png")
	defer imageFile.Close()

	validFileHolder := new(mockformfile.FileHolder)
	validFileHolder.On("FormFile", "logo").Return(imageFile, imageHeader, nil)

	noFileHolder := new(mockformfile.FileHolder)
	noFileHolder.On("FormFile", "logo").Return(nil, nil, http.ErrMissingFile)

	invalidFileHolder := new(mockformfile.FileHolder)
	invalidFileHolder.On("FormFile", "logo").Return(licenseFile, licenseHeader, nil)

	testCases := []testguard.InvalidParamsTestCase{
		{
			Description: "Should fail on missing ID",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "id",
			FileHolder:  validFileHolder,
			Sources: map[string]url.Values{
				"url": url.Values{},
			},
		},
		{
			Description: "Should fail on invalid uuid",
			MsgMatch:    params.ErrMsgInvalidUUID,
			FieldName:   "id",
			FileHolder:  validFileHolder,
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"not-a-uuid"},
				},
			},
		},
		{
			Description: "Should fail on missing logo",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "logo",
			FileHolder:  noFileHolder,
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"c3e98fdd-8a9e-4157-9a7c-fd2684e080ce"},
				},
			},
		},
		{
			Description: "Should fail on invalid logo",
			MsgMatch:    params.ErrMsgInvalidImage,
			FieldName:   "logo",
			FileHolder:  invalidFileHolder,
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"c3e98fdd-8a9e-4157-9a7c-fd2684e080ce"},
				},
			},
		},
	}

	// We wrap the tests otherwise the files will be closed too early
	// because they are all async
	t.Run("parallel wrapper", func(t *testing.T) {
		g := organizations.Endpoints[organizations.EndpointUploadLogo].Guard
		testguard.InvalidParams(t, g, testCases)
	})
}

func TestUploadLogoValidParams(t *testing.T) {
	cwd, _ := os.Getwd()

	testCases := []struct {
		description string
		sources     map[string]url.Values
		filename    string
	}{
		{
			"Should work with only a valid name",
			map[string]url.Values{
				"url": url.Values{
					"id": []string{"c3e98fdd-8a9e-4157-9a7c-fd2684e080ce"},
				},
			},
			"black_pixel.png",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()

			imageHeader, imageFile := testformfile.NewMultipartData(t, cwd, tc.filename)
			defer imageFile.Close()

			fileholder := new(mockformfile.FileHolder)
			fileholder.On("FormFile", "logo").Return(imageFile, imageHeader, nil)

			endpts := organizations.Endpoints[organizations.EndpointUploadLogo]
			data, err := endpts.Guard.ParseParams(tc.sources, fileholder)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*organizations.UploadLogoParams)
				assert.Equal(t, tc.sources["url"].Get("id"), p.ID)
			}
		})
	}
}

func TestUploadHappyPath(t *testing.T) {
	cwd, _ := os.Getwd()
	handlerParams := &organizations.UploadLogoParams{
		ID:   "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9",
		Logo: testformfile.NewFormFile(t, cwd, "black_pixel.png"),
	}
	defer handlerParams.Logo.File.Close()

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGet("*organizations.Organization", func(args mock.Arguments) {
		org := args.Get(0).(*organizations.Organization)
		org.ID = "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9"
		org.Name = "Google"
	})
	mockDB.ExpectUpdate("*organizations.Organization")

	// Mock the storage provider
	expectedURL := "http://domain.tld/image.png"
	storage := new(mockfilestorage.FileStorage)
	storage.ExpectWriteIfNotExist(false, expectedURL)
	storage.ExpectSetAttributes()

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*organizations.Payload", func(args mock.Arguments) {
		org := args.Get(0).(*organizations.Payload)
		assert.Equal(t, handlerParams.ID, org.ID, "ID should not have changed")
		assert.Equal(t, "Google", org.Name, "Name should not have changed")
		assert.NotNil(t, org.Logo, "Logo should not be nil")
		assert.Equal(t, expectedURL, org.Logo, "Logo should ha a URL set")
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)

	// call the handler
	deps := &router.Dependencies{DB: mockDB, Storage: storage}
	err := organizations.UploadLogo(req, deps)

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	storage.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestUploadNotFound(t *testing.T) {
	cwd, _ := os.Getwd()
	handlerParams := &organizations.UploadLogoParams{
		ID:   "0c2f0713-3f9b-4657-9cdd-2b4ed1f214e9",
		Logo: testformfile.NewFormFile(t, cwd, "black_pixel.png"),
	}
	defer handlerParams.Logo.File.Close()

	// Mock the database & add expectations
	mockDB := new(mockdb.DB)
	mockDB.ExpectGetNotFound("*organizations.Organization")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)

	// call the handler
	err := organizations.UploadLogo(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	e := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, e.HTTPStatus())
}
