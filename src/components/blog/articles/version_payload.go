package articles

import "github.com/Nivl/go-rest-tools/storage/db"

// VersionPayload represent an article Version that can be
// safely returned by the API
type VersionPayload struct {
	ID        string `json:"id"`
	ArticleID string `json:"article_id"`

	CreatedAt *db.Time `json:"created_at"`
	UpdatedAt *db.Time `json:"updated_at"`
	DeletedAt *db.Time `json:"deleted_at"`

	Title       string `json:"title"`
	Content     string `json:"content"`
	Subtitle    string `json:"subtitle"`
	Description string `json:"description"`
}

// VersionsPayload represent a list of Version that can be
// safely returned by the API
type VersionsPayload struct {
	Results []*VersionPayload `json:"results"`
}

// Export turns a Version into an object that is safe to be
// returned by the API
func (v *Version) Export() *VersionPayload {
	return &VersionPayload{
		ID:        v.ID,
		ArticleID: v.ArticleID,

		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
		DeletedAt: v.DeletedAt,

		Title:       v.Title,
		Content:     v.Content,
		Subtitle:    v.Subtitle,
		Description: v.Description,
	}
}

// Versions represents a list of version
type Versions []*Version

// Export turns a list of Version into an object that is safe to be
// returned by the API
func (vs Versions) Export() *VersionsPayload {
	pld := &VersionsPayload{}
	pld.Results = make([]*VersionPayload, len(vs))

	for i, v := range vs {
		pld.Results[i] = v.Export()
	}
	return pld
}
