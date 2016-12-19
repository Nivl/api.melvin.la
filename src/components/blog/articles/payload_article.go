package articles

import "github.com/melvin-laplanche/ml-api/src/db"
import "github.com/melvin-laplanche/ml-api/src/components/users"

// Payload represents an Article that can be safely returned by the API
type Payload struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`

	CreatedAt   db.Time              `json:"created_at"`
	UpdatedAt   db.Time              `json:"updated_at"`
	PublishedAt *db.Time             `json:"published_at,omitempty"`
	User        *users.PublicPayload `json:"user"`
	Content     *ContentPayload      `json:"content"`
	Draft       *ContentPayload      `json:"draft,omitempty"`
}

// Payloads is used to handle a list of Payload.
type Payloads struct {
	Results []*Payload `json:"results"`
}

// PublicExport turns an Article into an object that is safe to be
// returned by the API
func (a *Article) PublicExport() *Payload {
	return &Payload{
		ID:          a.ID,
		Slug:        a.Slug,
		CreatedAt:   *a.CreatedAt,
		UpdatedAt:   *a.UpdatedAt,
		PublishedAt: a.PublishedAt,
		User:        users.NewPublicPayload(a.User),
		Content:     a.Content.Export(),
	}
}

// PrivateExport turns an Article into an object that will contain private data
// that is safe to be returned by the API
func (a *Article) PrivateExport() *Payload {
	pld := a.PublicExport()
	pld.Draft = a.Draft.Export()
	return pld
}

// PublicExport turns a list of Articles into an object that is safe to be
// returned by the API
func (arts Articles) PublicExport() *Payloads {
	output := &Payloads{}

	output.Results = make([]*Payload, len(arts))
	for i, a := range arts {
		output.Results[i] = a.PublicExport()
	}
	return output
}

// PrivateExport turns a list of Articles into an object that is safe to be
// returned by the API, but contains privacy sensitive data
func (arts Articles) PrivateExport() *Payloads {
	output := &Payloads{}

	output.Results = make([]*Payload, len(arts))
	for i, a := range arts {
		output.Results[i] = a.PrivateExport()
	}
	return output
}
