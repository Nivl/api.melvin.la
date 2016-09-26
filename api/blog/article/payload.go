package article

import "github.com/Nivl/api.melvin.la/api/app/helpers"

// Exportable represents an Article that can be safely returned by the API
type Exportable struct {
	Title       string `json:"title"`
	Content     string `json:"content"`
	Slug        string `json:"slug"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

// NewPayloadFromModel turns an Article into an object that is safe to be
// returned by the API
func NewPayloadFromModel(a *Article) *Exportable {
	return &Exportable{
		Title:       a.Title,
		Content:     a.Content,
		Slug:        a.Slug,
		Subtitle:    a.Subtitle,
		Description: a.Description,
		CreatedAt:   helpers.GetDateForJSON(a.CreatedAt),
	}
}

// NewPayloadFromModels turns a []*Article into a list object that is safe to be
// returned by the API
func NewPayloadFromModels(list []*Article) []*Exportable {
	pld := make([]*Exportable, len(list))
	for i, a := range list {
		pld[i] = NewPayloadFromModel(a)
	}
	return pld
}
