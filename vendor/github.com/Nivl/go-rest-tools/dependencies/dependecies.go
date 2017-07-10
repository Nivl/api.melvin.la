package dependencies

import (
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/bsphere/le_go"
	"github.com/jmoiron/sqlx"
)

// DB represents an open connection with write access to the database
var DB db.DB

// InitPostgres inits the connection to the database
func InitPostgres(uri string) error {
	con, err := sqlx.Connect("postgres", uri)
	if err != nil {
		return err
	}

	// Unsafe returns a version of DB which will silently succeed to scan when
	// columns in the SQL result have no fields in the destination struct.
	DB = con.Unsafe()
	return nil
}

// Logentries represents an open connection to logentries
var Logentries *le_go.Logger

// InitLogentries inits the connection to logentries
func InitLogentries(token string) {
	le, err := le_go.Connect(token)
	if err != nil {
		panic(err)
	}
	Logentries = le
}

// Sendgrid is a sendgrid email client
var Sendgrid *SendgridParams

// SendgridParams represents the configuration of Sendgrid
type SendgridParams struct {
	APIKey         string
	From           string
	To             string
	StacktraceUUID string
}

// InitSendgrid creates a mailer that uses Sendgrid
func InitSendgrid(config *SendgridParams) {
	Sendgrid = config
}

// GCP represents the configuration of Google Cloud
type GCP struct {
	APIKey      string
	ProjectName string
	Bucket      string
}

// GoogleCloud contains the Google Cloud configuration
var GoogleCloud *GCP

// InitGCP setups Google Cloud Platform
func InitGCP(config *GCP) {
	GoogleCloud = config
}

// CloudinaryParams represents the configuration of Cloudinary
type CloudinaryParams struct {
	APIKey string
	Secret string
	Bucket string
}

// Cloudinary contains Cloudinary's configuration
var Cloudinary *CloudinaryParams

// InitCloudinary setups Cloudinary
func InitCloudinary(config *CloudinaryParams) {
	Cloudinary = config
}
