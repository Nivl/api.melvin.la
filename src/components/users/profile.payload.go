package users

import "github.com/Nivl/go-rest-tools/types/ptrs"

// ProfilePayload represents the public information of a user
type ProfilePayload struct {
	FirstName        string `json:"first_name,omitempty"`
	LastName         string `json:"last_name,omitempty"`
	Picture          string `json:"picture,omitempty"`
	PhoneNumber      string `json:"phone_number,omitempty"`
	PublicEmail      string `json:"public_email,omitempty"`
	LinkedIn         string `json:"linkedin_custom_url,omitempty"`
	FacebookUsername string `json:"facebook_username,omitempty"`
	TwitterUsername  string `json:"twitter_username,omitempty"`

	*Payload // User payload
}

// ExportPublic returns a ProfilePayload containing only the fields that are safe to
// be seen by anyone
func (p *Profile) ExportPublic() *ProfilePayload {
	// It's OK to export a nil experience
	if p == nil {
		return nil
	}

	return &ProfilePayload{
		FirstName:        ptrs.UnwrapString(p.FirstName),
		LastName:         ptrs.UnwrapString(p.LastName),
		Picture:          ptrs.UnwrapString(p.Picture),
		PhoneNumber:      ptrs.UnwrapString(p.PhoneNumber),
		PublicEmail:      ptrs.UnwrapString(p.PublicEmail),
		LinkedIn:         ptrs.UnwrapString(p.LinkedIn),
		FacebookUsername: ptrs.UnwrapString(p.FacebookUsername),
		TwitterUsername:  ptrs.UnwrapString(p.TwitterUsername),
		Payload:          NewPayload(p.User),
	}
}

// ExportPrivate returns a ProfilePayload containing all the fields
func (p *Profile) ExportPrivate() *ProfilePayload {
	// It's OK to export a nil experience
	if p == nil {
		return nil
	}

	pld := p.ExportPublic()
	pld.Payload = NewPrivatePayload(p.User)
	return pld
}