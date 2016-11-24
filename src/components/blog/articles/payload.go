package articles

import "github.com/melvin-laplanche/ml-api/src/db"
import "github.com/melvin-laplanche/ml-api/src/components/users"

// PublicPayload represents an Article that can be safely returned by the API
type PublicPayload struct {
	ID          string               `json:"id"`
	Title       string               `json:"title"`
	Content     string               `json:"content"`
	Slug        string               `json:"slug"`
	Subtitle    string               `json:"subtitle"`
	Description string               `json:"description"`
	CreatedAt   db.Time              `json:"created_at"`
	UpdatedAt   db.Time              `json:"updated_at"`
	IsPublished bool                 `json:"is_published"`
	User        *users.PublicPayload `json:"user"`
}

// PublicPayloads is used to handle a list of publicPayload.
type PublicPayloads struct {
	Results []*PublicPayload `json:"results"`
}

// Export turns an Article into an object that is safe to be
// returned by the API
func (a *Article) Export() *PublicPayload {
	return &PublicPayload{
		ID:          a.ID,
		Title:       a.Title,
		Content:     a.Content,
		Slug:        a.Slug,
		Subtitle:    a.Subtitle,
		Description: a.Description,
		CreatedAt:   *a.CreatedAt,
		UpdatedAt:   *a.UpdatedAt,
		IsPublished: a.IsPublished,
		User:        users.NewPublicPayload(&a.User),
	}
}

// Export turns a list of Articles into an object that is safe to be
// returned by the API
func (arts Articles) Export() *PublicPayloads {
	output := &PublicPayloads{}

	output.Results = make([]*PublicPayload, len(arts))
	for i, a := range arts {
		output.Results[i] = a.Export()
	}
	return output
}
