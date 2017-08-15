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

// Setup parses the env, sets the app globals and returns the params
func Setup() (*Args, error) {
	var params Args
	if err := envconfig.Process("", &params); err != nil {
		return nil, err
	}

	if err := dependencies.InitPostgres(params.PostgresURI); err != nil {
		return nil, err
	}

	if params.LogEntriesToken != "" {
		dependencies.InitLogentries(params.LogEntriesToken)
	}

	if params.EmailAPIKey != "" {
		p := &dependencies.SendgridParams{
			APIKey:         params.EmailAPIKey,
			From:           params.EmailFrom,
			To:             params.EmailTo,
			StacktraceUUID: params.EmailStacktraceUUID,
		}
		dependencies.InitSendgrid(p)
	}

	if params.GCPAPIKey != "" {
		p := &dependencies.GCP{
			APIKey:      params.GCPAPIKey,
			ProjectName: params.GCPProject,
			Bucket:      params.GCPBucket,
		}
		dependencies.InitGCP(p)
	}

	if params.CloudinaryAPIKey != "" {
		p := &dependencies.CloudinaryParams{
			APIKey: params.CloudinaryAPIKey,
			Secret: params.CloudinarySecret,
			Bucket: params.CloudinaryBucket,
		}
		dependencies.InitCloudinary(p)
	}

	return &params, nil
}
