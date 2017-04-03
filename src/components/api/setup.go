package api

import (
	"github.com/Nivl/go-rest-tools/logger"
	"github.com/Nivl/go-rest-tools/notifiers/mailer"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/bsphere/le_go"
	"github.com/kelseyhightower/envconfig"
)

// Args represents the app args
type Args struct {
	Port            string `default:"5000"`
	PostgresURI     string `required:"true" envconfig:"postgres_uri"`
	LogEntriesToken string `envconfig:"logentries_token"`
	EmailAPIKey     string `envconfig:"email_api_key"`
	EmailFrom       string `envconfig:"email_default_from"`
	EmailTo         string `envconfig:"email_default_to"`
	Debug           bool   `default:"false"`
}

// Setup parses the env, sets the app globals and returns the params
func Setup() *Args {
	var params Args
	if err := envconfig.Process("", &params); err != nil {
		panic(err)
	}

	if err := db.Setup(params.PostgresURI); err != nil {
		panic(err)
	}

	// LogEntries
	if params.LogEntriesToken != "" {
		le, err := le_go.Connect(params.LogEntriesToken)
		if err != nil {
			panic(err)
		}
		logger.LogEntries = le
	}

	// Sendgrid
	if params.EmailAPIKey != "" {
		mailer.Emailer = mailer.NewMailer(params.EmailAPIKey, params.EmailFrom, params.EmailTo)
	}

	return &params
}

// SetupIfNeeded parses the env, and sets the app globals
func SetupIfNeeded() {
	if db.Writer != nil {
		return
	}

	Setup()
}
