package sessions

import "github.com/Nivl/api.melvin.la/api/auth"

// Payload represents a Session that can be safely returned by the API
type Payload struct {
	Token  string `json:"token"`
	UserID string `json:"user_id"`
}

// NewPayloadFromModel turns a Session into an object that is safe to be
// returned by the API
func NewPayloadFromModel(s *auth.Session) *Payload {
	return &Payload{
		Token:  s.ID.Hex(),
		UserID: s.UserID.Hex(),
	}
}
