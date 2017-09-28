package reporter

import "github.com/Nivl/go-rest-tools/security/auth"

// Creator describes a struct used to create reporters
type Creator interface {
	New() (Reporter, error)
}

// Reporter describes a struct used to report errors
type Reporter interface {
	SetUser(u *auth.User)
	AddTag(key, value string)
	AddTags(tags map[string]string)
	ReportError(err error)
	ReportErrorAndWait(err error)
}
