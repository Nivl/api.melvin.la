package users

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
		FirstName:        p.FirstName,
		LastName:         p.LastName,
		Picture:          p.Picture,
		PhoneNumber:      p.PhoneNumber,
		PublicEmail:      p.PublicEmail,
		LinkedIn:         p.LinkedIn,
		FacebookUsername: p.FacebookUsername,
		TwitterUsername:  p.TwitterUsername,
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
