package reporter

import (
	"github.com/Nivl/go-rest-tools/security/auth"
)

var (
	// Makes sure the Email object implement Reporter
	_ Reporter = (*Noop)(nil)
	// Makes sure NoopCreator is a Creator
	_ Creator = (*NoopCreator)(nil)
)

// NewNoopCreator creates a client's creator
func NewNoopCreator() (*NoopCreator, error) {
	return &NoopCreator{}, nil
}

// NoopCreator is a Creator used to get no-op reporter
type NoopCreator struct{}

// New returns a new mailer Reporter
func (c *NoopCreator) New() (Reporter, error) {
	return NewNoop()
}

// NewNoop creates a new Mailer reporter
func NewNoop() (*Noop, error) {
	return &Noop{}, nil
}

// Noop represents reporter that does nothin
type Noop struct {
}

// SetUser does nothing
func (r *Noop) SetUser(u *auth.User) {}

// AddTag does nothing
func (r *Noop) AddTag(key, value string) {}

// AddTags does nothing
func (r *Noop) AddTags(tags map[string]string) {}

// ReportError does nothing
func (r *Noop) ReportError(err error) {}

// CaptureErrorAndWait does nothing
func (r *Noop) ReportErrorAndWait(err error) {}
