package reporter

import (
	"runtime/debug"

	"github.com/Nivl/go-rest-tools/notifiers/mailer"
	"github.com/Nivl/go-rest-tools/security/auth"
)

var (
	// Makes sure the Email object implement Reporter
	_ Reporter = (*Mailer)(nil)
	// Makes sure MailerCreator is a Creator
	_ Creator = (*MailerCreator)(nil)
)

// NewMailerCreator creates a client's creator
func NewMailerCreator(m mailer.Mailer) (*MailerCreator, error) {
	return &MailerCreator{
		mailer: m,
	}, nil
}

// MailerCreator is a Creator used to get mailers
type MailerCreator struct {
	mailer mailer.Mailer
}

// New returns a new mailer Reporter
func (c *MailerCreator) New() (Reporter, error) {
	return NewMailer(c.mailer)
}

// NewMailer creates a new Mailer reporter
func NewMailer(m mailer.Mailer) (*Mailer, error) {
	return &Mailer{
		client: m,
		tags:   map[string]string{},
	}, nil
}

// Mailer represents a client used to report errors by email
type Mailer struct {
	client mailer.Mailer
	user   *auth.User
	tags   map[string]string
}

// SetUser attaches the provided user to report
func (r *Mailer) SetUser(u *auth.User) {
	r.user = u
}

// AddTag attaches the provided data to the report
func (r *Mailer) AddTag(key, value string) {
	r.tags[key] = value
}

// AddTags attaches the provided data to the report
func (r *Mailer) AddTags(tags map[string]string) {
	for key, value := range tags {
		r.tags[key] = value
	}
}

// CaptureError sends a report with this specific error
func (r *Mailer) ReportError(err error) {
	// We copy the tags so we can add more data without affecting the source
	finalTags := map[string]string{}
	if r.user != nil {
		finalTags["User ID"] = r.user.ID
		finalTags["User email"] = r.user.Email
		finalTags["User name"] = r.user.Name
	}
	for k, v := range r.tags {
		finalTags[k] = v
	}

	sendEmail := func(stacktrace []byte) {
		r.client.SendStackTrace(stacktrace, err.Error(), finalTags)
	}
	go sendEmail(debug.Stack())
}

// ReportErrorAndWait sends a report with this specific error
func (r *Mailer) ReportErrorAndWait(err error) {
	// We copy the tags so we can add more data without affecting the source
	finalTags := map[string]string{}
	if r.user != nil {
		finalTags["User ID"] = r.user.ID
		finalTags["User email"] = r.user.Email
		finalTags["User name"] = r.user.Name
	}
	for k, v := range r.tags {
		finalTags[k] = v
	}

	r.client.SendStackTrace(debug.Stack(), err.Error(), finalTags)
}
