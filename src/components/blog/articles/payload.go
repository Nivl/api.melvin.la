package articles

import "github.com/melvin-laplanche/ml-api/src/app/helpers"

// PublicPayload represents an Article that can be safely returned by the API
type PublicPayload struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Slug        string `json:"slug"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
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
		CreatedAt:   helpers.GetDateForJSON(a.CreatedAt),
		UpdatedAt:   helpers.GetDateForJSON(a.UpdatedAt),
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
