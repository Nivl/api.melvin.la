package articles

import "github.com/melvin-laplanche/ml-api/src/db"
import "github.com/melvin-laplanche/ml-api/src/components/users"

// PublicPayload represents an Article that can be safely returned by the API
type PublicPayload struct {
	ID   string `json:"id"`
	Slug string `json:"slug"`

	CreatedAt   db.Time              `json:"created_at"`
	UpdatedAt   db.Time              `json:"updated_at"`
	PublishedAt *db.Time             `json:"is_published,omitempty"`
	User        *users.PublicPayload `json:"user"`
	Content     *ContentPayload      `json:"content"`
}

// PublicPayloads is used to handle a list of publicPayload.
type PublicPayloads struct {
	Results []*PublicPayload `json:"results"`
}

// PrivatePayload represents an Article containing sensitive data that can returned by the API
type PrivatePayload struct {
	PublicPayload
	Draft *ContentPayload `json:"draft"`
}

// PublicExport turns an Article into an object that is safe to be
// returned by the API
func (a *Article) PublicExport() *PublicPayload {
	return &PublicPayload{
		ID:          a.ID,
		Slug:        a.Slug,
		CreatedAt:   *a.CreatedAt,
		UpdatedAt:   *a.UpdatedAt,
		PublishedAt: a.PublishedAt,
		User:        users.NewPublicPayload(&a.User),
		Content:     a.Content.Export(),
	}
}

// PrivateExport turns an Article into an object that will contain private data
// that is safe to be returned by the API
func (a *Article) PrivateExport() *PrivatePayload {
	return &PrivatePayload{
		PublicPayload: *a.PublicExport(),
		Draft:         a.Draft.Export(),
	}
}

// Export turns a list of Articles into an object that is safe to be
// returned by the API
func (arts Articles) Export() *PublicPayloads {
	output := &PublicPayloads{}

	output.Results = make([]*PublicPayload, len(arts))
	for i, a := range arts {
		output.Results[i] = a.PublicExport()
	}
	return output
}
