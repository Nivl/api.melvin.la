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
	Debug               bool   `default:"false"`
}

// Setup parses the env, sets the app globals and returns the params
func Setup() *Args {
	var params Args
	if err := envconfig.Process("", &params); err != nil {
		panic(err)
	}

	if err := dependencies.InitPostgres(params.PostgresURI); err != nil {
		panic(err)
	}

	if params.LogEntriesToken != "" {
		dependencies.InitLogentries(params.LogEntriesToken)
	}

	if params.EmailAPIKey != "" {
		dependencies.InitSendgrid(params.EmailAPIKey, params.EmailFrom, params.EmailTo, params.EmailStacktraceUUID)
	}

	return &params
}
