package reporter

import (
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/getsentry/raven-go"
)

var (
	// Makes sure the Sentry object implement Reporter
	_ Reporter = (*Sentry)(nil)
	// Makes sure SentryCreator is a Creator
	_ Creator = (*SentryCreator)(nil)
)

// NewSentryCreator creates a client's creator
func NewSentryCreator(con string) (*SentryCreator, error) {
	return &SentryCreator{
		con: con,
	}, nil
}

// SentryCreator is a Creator use to get sentry clients
type SentryCreator struct {
	con string
}

// New returns a new sentry client
func (c *SentryCreator) New() (Reporter, error) {
	return NewSentry(c.con)
}

// NewSentry creates a new sentry client
func NewSentry(con string) (*Sentry, error) {
	c, err := raven.New(con)
	if err != nil {
		return nil, err
	}

	return &Sentry{
		client: c,
		tags:   map[string]string{},
	}, nil
}

// Sentry represents a client used to report errors to Sentry
type Sentry struct {
	client *raven.Client
	tags   map[string]string
}

// SetUser attaches the provided user to report
func (r *Sentry) SetUser(u *auth.User) {
	sentryUser := &raven.User{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Name,
	}
	r.client.SetUserContext(sentryUser)
}

// AddTag attaches the provided data to the report
func (r *Sentry) AddTag(key, value string) {
	r.tags[key] = value
}

// AddTags attaches the provided data to the report
func (r *Sentry) AddTags(tags map[string]string) {
	for key, value := range tags {
		r.tags[key] = value
	}
}

// ReportError sends a report with the specified error
func (r *Sentry) ReportError(err error) {
	r.client.CaptureError(err, r.tags)
}

// ReportErrorAndWait sends a report with the specified error, and wait for
// the reports to be sent
func (r *Sentry) ReportErrorAndWait(err error) {
	r.client.CaptureErrorAndWait(err, r.tags)
}
