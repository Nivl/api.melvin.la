package users_test

import (
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/satori/go.uuid"

	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/router/formfile/mockformfile"
	"github.com/Nivl/go-rest-tools/router/formfile/testformfile"
	"github.com/Nivl/go-rest-tools/router/guard/testguard"
	"github.com/Nivl/go-rest-tools/router/mockrouter"
	"github.com/Nivl/go-rest-tools/router/params"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db/mockdb"
	"github.com/Nivl/go-rest-tools/storage/filestorage/mockfilestorage"
	"github.com/Nivl/go-rest-tools/types/apierror"
	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/melvin-laplanche/ml-api/src/components/users/testusers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadPictureAccess(t *testing.T) {
	testCases := []testguard.AccessTestCase{
		{
			Description: "Should fail for anonymous users",
			User:        nil,
			ErrCode:     http.StatusUnauthorized,
		},
		{
			Description: "Should work for logged users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0"},
			ErrCode:     0,
		},
		{
			Description: "Should work for admin users",
			User:        &auth.User{ID: "48d0c8b8-d7a3-4855-9d90-29a06ef474b0", IsAdmin: true},
			ErrCode:     0,
		},
	}

	g := users.Endpoints[users.EndpointUploadPicture].Guard
	testguard.AccessTest(t, g, testCases)
}

func TestUploadPictureInvalidParams(t *testing.T) {

	// create the multipart data
	cwd, _ := os.Getwd()
	licenseHeader, licenseFile := testformfile.NewMultipartData(t, cwd, "LICENSE")
	defer licenseFile.Close()

	imageHeader, imageFile := testformfile.NewMultipartData(t, cwd, "black_pixel.png")
	defer imageFile.Close()

	validFileHolder := new(mockformfile.FileHolder)
	validFileHolder.On("FormFile", "picture").Return(imageFile, imageHeader, nil)

	noFileHolder := new(mockformfile.FileHolder)
	noFileHolder.On("FormFile", "picture").Return(nil, nil, http.ErrMissingFile)

	invalidFileHolder := new(mockformfile.FileHolder)
	invalidFileHolder.On("FormFile", "picture").Return(licenseFile, licenseHeader, nil)

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
			Description: "Should fail on missing picture",
			MsgMatch:    params.ErrMsgMissingParameter,
			FieldName:   "picture",
			FileHolder:  noFileHolder,
			Sources: map[string]url.Values{
				"url": url.Values{
					"id": []string{"c3e98fdd-8a9e-4157-9a7c-fd2684e080ce"},
				},
			},
		},
		{
			Description: "Should fail on invalid picture",
			MsgMatch:    params.ErrMsgInvalidImage,
			FieldName:   "picture",
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
		g := users.Endpoints[users.EndpointUploadPicture].Guard
		testguard.InvalidParams(t, g, testCases)
	})
}

func TestUploadPictureValidParams(t *testing.T) {
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
			fileholder.On("FormFile", "picture").Return(imageFile, imageHeader, nil)

			endpts := users.Endpoints[users.EndpointUploadPicture]
			data, err := endpts.Guard.ParseParams(tc.sources, fileholder)
			assert.NoError(t, err)

			if data != nil {
				p := data.(*users.UploadPictureParams)
				assert.Equal(t, tc.sources["url"].Get("id"), p.ID)
			}
		})
	}
}

func TestUploadHappyPath(t *testing.T) {
	profile := testusers.NewProfile()

	cwd, _ := os.Getwd()
	handlerParams := &users.UploadPictureParams{
		ID:      profile.User.ID,
		Picture: testformfile.NewFormFile(t, cwd, "black_pixel.png"),
	}
	defer handlerParams.Picture.File.Close()

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGet("*users.Profile", func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *profile
	})
	mockDB.ExpectUpdate("*users.Profile")

	// Mock the storage provider
	expectedURL := "http://domain.tld/image.png"
	storage := new(mockfilestorage.FileStorage)
	storage.ExpectWriteIfNotExist(false, expectedURL)
	storage.ExpectSetAttributes()

	// Mock the response & add expectations
	res := new(mockrouter.HTTPResponse)
	res.ExpectOk("*users.ProfilePayload", func(args mock.Arguments) {
		p := args.Get(0).(*users.ProfilePayload)
		assert.NotNil(t, p.Picture, "Picture should not be nil")
		assert.Equal(t, expectedURL, p.Picture, "Picture should ha a URL set")
	})

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Response").Return(res)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(profile.User)

	// call the handler
	deps := &router.Dependencies{DB: mockDB, Storage: storage}
	err := users.UploadPicture(req, deps)

	// Assert everything
	assert.NoError(t, err, "the handler should not have fail")
	storage.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)
	res.AssertExpectations(t)
}

func TestUploadAdminNotFoundUser(t *testing.T) {
	profile := testusers.NewProfile()
	profile.User.IsAdmin = true

	cwd, _ := os.Getwd()
	handlerParams := &users.UploadPictureParams{
		ID:      uuid.NewV4().String(),
		Picture: testformfile.NewFormFile(t, cwd, "black_pixel.png"),
	}
	defer handlerParams.Picture.File.Close()

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGetNotFound("*users.Profile")

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(profile.User)

	// call the handler
	err := users.UploadPicture(req, &router.Dependencies{DB: mockDB})

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	e := apierror.Convert(err)
	assert.Equal(t, http.StatusNotFound, e.HTTPStatus())
}

func TestUploadStorageFailed(t *testing.T) {
	profile := testusers.NewProfile()

	cwd, _ := os.Getwd()
	handlerParams := &users.UploadPictureParams{
		ID:      profile.User.ID,
		Picture: testformfile.NewFormFile(t, cwd, "black_pixel.png"),
	}
	defer handlerParams.Picture.File.Close()

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGet("*users.Profile", func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *(testusers.NewProfile())
	})

	// Mock the storage provider
	storage := new(mockfilestorage.FileStorage)
	storage.ExpectWriteIfNotExistError()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(profile.User)

	// call the handler
	deps := &router.Dependencies{DB: mockDB, Storage: storage}
	err := users.UploadPicture(req, deps)

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	storage.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}

func TestUploadDBNoCon(t *testing.T) {
	profile := testusers.NewProfile()

	cwd, _ := os.Getwd()
	handlerParams := &users.UploadPictureParams{
		ID:      profile.User.ID,
		Picture: testformfile.NewFormFile(t, cwd, "black_pixel.png"),
	}
	defer handlerParams.Picture.File.Close()

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}
	mockDB.ExpectGet("*users.Profile", func(args mock.Arguments) {
		p := args.Get(0).(*users.Profile)
		*p = *(testusers.NewProfile())
	})
	mockDB.ExpectUpdateError("*users.Profile")

	// Mock the storage provider
	expectedURL := "http://domain.tld/image.png"
	storage := new(mockfilestorage.FileStorage)
	storage.ExpectWriteIfNotExist(false, expectedURL)
	storage.ExpectSetAttributes()

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(profile.User)

	// call the handler
	deps := &router.Dependencies{DB: mockDB, Storage: storage}
	err := users.UploadPicture(req, deps)

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	storage.AssertExpectations(t)
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusInternalServerError, apiError.HTTPStatus())
}

func TestUploadWrongUser(t *testing.T) {
	profile := testusers.NewProfile()

	cwd, _ := os.Getwd()
	handlerParams := &users.UploadPictureParams{
		ID:      uuid.NewV4().String(),
		Picture: testformfile.NewFormFile(t, cwd, "black_pixel.png"),
	}
	defer handlerParams.Picture.File.Close()

	// Mock the database & add expectations
	mockDB := &mockdb.Connection{}

	// Mock the request & add expectations
	req := new(mockrouter.HTTPRequest)
	req.On("Params").Return(handlerParams)
	req.On("User").Return(profile.User)

	// call the handler
	deps := &router.Dependencies{DB: mockDB}
	err := users.UploadPicture(req, deps)

	// Assert everything
	assert.Error(t, err, "the handler should have fail")
	mockDB.AssertExpectations(t)
	req.AssertExpectations(t)

	apiError := apierror.Convert(err)
	assert.Equal(t, http.StatusForbidden, apiError.HTTPStatus())
}
