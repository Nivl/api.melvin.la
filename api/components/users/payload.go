package users

// Exportable represents an Article that can be safely returned by the API
import "github.com/Nivl/api.melvin.la/api/auth"

// PrivatePayload represents a user payload with non public field
type PrivatePayload struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// NewPrivatePayload turns a user into an object that is safe to be
// returned by the API
func NewPrivatePayload(u *auth.User) *PrivatePayload {
	return &PrivatePayload{
		ID:    u.ID.Hex(),
		Name:  u.Name,
		Email: u.Email,
	}
}
