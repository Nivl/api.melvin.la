package articles

import (
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/melvin-laplanche/ml-api/src/components/users"
)

// Payload represents an Article that can be safely returned by the API
type Payload struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`

	Title       string `json:"title"`
	Content     string `json:"content"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`

	CreatedAt   *db.Time `json:"created_at"`
	UpdatedAt   *db.Time `json:"updated_at"`
	PublishedAt *db.Time `json:"published_at,omitempty"`

	User *users.Payload `json:"user"`
}

// Payloads is used to handle a list of Payload.
type Payloads struct {
	Results []*Payload `json:"results"`
}

// Export turns an Article into an object that is safe to be
// returned by the API
func (a *Article) Export() *Payload {
	return &Payload{
		ID:          a.ID,
		Slug:        a.Slug,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
		PublishedAt: a.PublishedAt,
		User:        users.NewPayload(a.User),
		Title:       a.Version.Title,
		Content:     a.Version.Content,
		Subtitle:    a.Version.Subtitle,
		Description: a.Version.Description,
	}
}

// Export turns a list of Articles into an object that is safe to be
// returned by the API
func (arts Articles) Export() *Payloads {
	output := &Payloads{}

	output.Results = make([]*Payload, len(arts))
	for i, a := range arts {
		output.Results[i] = a.Export()
	}
	return output
}
