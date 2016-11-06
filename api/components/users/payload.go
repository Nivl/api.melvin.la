package users

import "github.com/Nivl/api.melvin.la/api/auth"

// PrivatePayload represents a user payload with non public field
type PrivatePayload struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// NewPrivatePayload turns a user into an object that is safe to be
// returned by the API
func NewPrivatePayload(u *auth.User) *PrivatePayload {
	return &PrivatePayload{
		UUID:  u.UUID,
		Name:  u.Name,
		Email: u.Email,
	}
}

// PublicPayload represents a user payload with no private field
type PublicPayload struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
}

// NewPublicPayload turns a user into an object that is safe to be
// returned by the API
func NewPublicPayload(u *auth.User) *PublicPayload {
	return &PublicPayload{
		UUID: u.UUID,
		Name: u.Name,
	}
}
