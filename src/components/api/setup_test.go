package api_test

import (
	"errors"
	"testing"

	"github.com/Nivl/go-rest-tools/dependencies/mockdependencies"
	"github.com/melvin-laplanche/ml-api/src/components/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// stringType represents a string argument
var stringType = mock.AnythingOfType("string")

func TestSetupRequired(t *testing.T) {
	args := &api.Args{}
	deps := &mockdependencies.Dependencies{}

	args.PostgresURI = "postgres uri"
	deps.On("SetDB", args.PostgresURI).Return(nil)

	assert.NoError(t, api.Setup(args, deps), "Setup should have worked")
	deps.AssertExpectations(t)
}

func TestSetupHappyPath(t *testing.T) {
	args := &api.Args{}
	deps := &mockdependencies.Dependencies{}

	args.PostgresURI = "postgres uri"
	deps.On("SetDB", args.PostgresURI).Return(nil)

	args.LogEntriesToken = "logentry token"
	deps.On("SetLogentries", args.LogEntriesToken).Return(nil)

	args.EmailAPIKey = "email api key"
	args.EmailFrom = "email from"
	args.EmailTo = "email to"
	args.EmailStacktraceUUID = "email stacktrace uuid"
	deps.On("SetSendgrid", args.EmailAPIKey, args.EmailFrom, args.EmailTo, args.EmailStacktraceUUID).Return(nil)

	args.GCPAPIKey = "gcp api key"
	args.GCPProject = "gcp project"
	args.GCPBucket = "gcp bucket"
	deps.On("SetGCP", args.GCPAPIKey, args.GCPProject, args.GCPBucket).Return(nil)

	args.CloudinaryAPIKey = "cloudinary api key"
	args.CloudinarySecret = "cloudinary project"
	args.CloudinaryBucket = "cloudinary bucket"
	deps.On("SetCloudinary", args.CloudinaryAPIKey, args.CloudinarySecret, args.CloudinaryBucket).Return(nil)

	args.SentryDSN = "sentry dsn"
	deps.On("SetSentry", args.SentryDSN).Return(nil)

	assert.NoError(t, api.Setup(args, deps), "Setup should have worked")
	deps.AssertExpectations(t)
}

func TestSetupErrorDB(t *testing.T) {
	args := &api.Args{}
	deps := &mockdependencies.Dependencies{}

	errWanted := errors.New("db error")
	deps.On("SetDB", stringType).Return(errWanted)

	err := api.Setup(args, deps)
	assert.Error(t, err, "Setup should have failed")
	assert.Equal(t, errWanted, err, "Wrong error returned")
	deps.AssertExpectations(t)
}

func TestSetupErrorLogentries(t *testing.T) {
	args := &api.Args{}
	deps := &mockdependencies.Dependencies{}
	deps.On("SetDB", stringType).Return(nil)

	args.LogEntriesToken = "token"
	errWanted := errors.New("logentries error")
	deps.On("SetLogentries", args.LogEntriesToken).Return(errWanted)

	err := api.Setup(args, deps)
	assert.Error(t, err, "Setup should have failed")
	assert.Equal(t, errWanted, err, "Wrong error returned")
	deps.AssertExpectations(t)
}

func TestSetupErrorSendgrid(t *testing.T) {
	args := &api.Args{}
	deps := &mockdependencies.Dependencies{}
	deps.On("SetDB", stringType).Return(nil)

	args.EmailAPIKey = "token"
	errWanted := errors.New("sendgrid error")
	deps.On("SetSendgrid", args.EmailAPIKey, stringType, stringType, stringType).Return(errWanted)

	err := api.Setup(args, deps)
	assert.Error(t, err, "Setup should have failed")
	assert.Equal(t, errWanted, err, "Wrong error returned")
	deps.AssertExpectations(t)
}

func TestSetupErrorGCP(t *testing.T) {
	args := &api.Args{}
	deps := &mockdependencies.Dependencies{}
	deps.On("SetDB", stringType).Return(nil)

	args.GCPAPIKey = "token"
	errWanted := errors.New("gcp error")
	deps.On("SetGCP", args.GCPAPIKey, stringType, stringType).Return(errWanted)

	err := api.Setup(args, deps)
	assert.Error(t, err, "Setup should have failed")
	assert.Equal(t, errWanted, err, "Wrong error returned")
	deps.AssertExpectations(t)
}

func TestSetupErrorCloudinary(t *testing.T) {
	args := &api.Args{}
	deps := &mockdependencies.Dependencies{}
	deps.On("SetDB", stringType).Return(nil)

	args.CloudinaryAPIKey = "token"
	errWanted := errors.New("cloudinary error")
	deps.On("SetCloudinary", args.CloudinaryAPIKey, stringType, stringType).Return(errWanted)

	err := api.Setup(args, deps)
	assert.Error(t, err, "Setup should have failed")
	assert.Equal(t, errWanted, err, "Wrong error returned")
	deps.AssertExpectations(t)
}

func TestSetupErrorSentry(t *testing.T) {
	args := &api.Args{}
	deps := &mockdependencies.Dependencies{}
	deps.On("SetDB", stringType).Return(nil)

	errWanted := errors.New("sentry error")
	args.SentryDSN = "sentry dsn"
	deps.On("SetSentry", args.SentryDSN).Return(errWanted)

	err := api.Setup(args, deps)
	assert.Error(t, err, "Setup should have failed")
	assert.Equal(t, errWanted, err, "Wrong error returned")
	deps.AssertExpectations(t)
}
