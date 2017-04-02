package articles

// // ContentPayload represent an article content (or draft) that can be
// // safely returned by the API
// type ContentPayload struct {
// 	Title       string `json:"title"`
// 	Content     string `json:"content"`
// 	Subtitle    string `json:"subtitle"`
// 	Description string `json:"description"`
// }

// // Export turns an Article into an object that is safe to be
// // returned by the API
// func (c *Content) Export() *ContentPayload {
// 	return &ContentPayload{
// 		Title:       c.Title,
// 		Content:     c.Content,
// 		Subtitle:    c.Subtitle,
// 		Description: c.Description,
// 	}
// }

// // Export turns an Article into an object that is safe to be
// // returned by the API
// func (d *Draft) Export() *ContentPayload {
// 	if d == nil {
// 		return nil
// 	}

// 	// A draft is just a content
// 	content := Content(*d)
// 	return (&content).Export()
// }
