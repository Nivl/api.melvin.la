package sessions

import "github.com/Nivl/api.melvin.la/api/auth"

// Payload represents a Session that can be safely returned by the API
type Payload struct {
	Token    string `json:"token"`
	UserUUID string `json:"user_uuid"`
}

// NewPayloadFromModel turns a Session into an object that is safe to be
// returned by the API
func NewPayloadFromModel(s *auth.Session) *Payload {
	return &Payload{
		Token:    s.UUID,
		UserUUID: s.UserUUID,
	}
}
