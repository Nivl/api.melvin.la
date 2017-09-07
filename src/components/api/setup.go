package api

import (
	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/kelseyhightower/envconfig"
)

// Args represents the app args
type Args struct {
	Port                string `default:"5000"`
	PostgresURI         string `required:"true" envconfig:"postgres_uri"`
	LogEntriesToken     string `envconfig:"logentries_token"`
	EmailAPIKey         string `envconfig:"email_api_key"`
	EmailFrom           string `envconfig:"email_default_from"`
	EmailTo             string `envconfig:"email_default_to"`
	EmailStacktraceUUID string `envconfig:"email_stacktrace_uuid"`
	CloudinaryAPIKey    string `envconfig:"cloudinary_api_key"`
	CloudinarySecret    string `envconfig:"cloudinary_secret"`
	CloudinaryBucket    string `envconfig:"cloudinary_bucket"`
	GCPAPIKey           string `envconfig:"gcp_api_key"`
	GCPProject          string `envconfig:"gcp_project"`
	GCPBucket           string `envconfig:"gcp_bucket"`
	Debug               bool   `default:"false"`
}

// DefaultSetup parses the env and returns the args and dependencies
func DefaultSetup() (*Args, dependencies.Dependencies, error) {
	params := &Args{}
	if err := envconfig.Process("", params); err != nil {
		return nil, nil, err
	}

	deps := &dependencies.APIDependencies{}
	err := Setup(params, deps)
	return params, deps, err
}

// Setup parses the env, sets the app globals and returns the params
func Setup(params *Args, deps dependencies.Dependencies) error {
	if err := deps.SetDB(params.PostgresURI); err != nil {
		return err
	}

	if params.LogEntriesToken != "" {
		if err := deps.SetLogentries(params.LogEntriesToken); err != nil {
			return err
		}
	}

	if params.EmailAPIKey != "" {
		if err := deps.SetSendgrid(params.EmailAPIKey, params.EmailFrom, params.EmailTo, params.EmailStacktraceUUID); err != nil {
			return err
		}
	}

	if params.GCPAPIKey != "" {
		if err := deps.SetGCP(params.GCPAPIKey, params.GCPProject, params.GCPBucket); err != nil {
			return err
		}
	}

	if params.CloudinaryAPIKey != "" {
		if err := deps.SetCloudinary(params.CloudinaryAPIKey, params.CloudinarySecret, params.CloudinaryBucket); err != nil {
			return err
		}
	}

	return nil
}
