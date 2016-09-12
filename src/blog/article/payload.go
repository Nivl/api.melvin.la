package article

import "github.com/Nivl/api.melvin.la/src/app/helpers"

// Payload represents an Article that can be sent to or returned by the API
type Payload struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Slug        string `json:"slug"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

// NewPayloadFromModel turns an Article into a an object that is safe to be
// returned by the API
func NewPayloadFromModel(a *Article) *Payload {
	return &Payload{
		Title:       a.Title,
		Content:     a.Content,
		Slug:        a.Slug,
		Subtitle:    a.Subtitle,
		Description: a.Description,
		CreatedAt:   helpers.GetDateForJSON(a.CreatedAt),
	}
}
