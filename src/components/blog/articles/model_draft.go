package articles

// Draft Represents an article draft
type Draft Content

// ToContent returns a Content from a Draft
func (d *Draft) ToContent() *Content {
	content := Content(*d)
	return &content
}

// Create persists a draft in the database
func (d *Draft) Create() error {
	return d.ToContent().Create()
}

// Update updates most of the fields of a persisted content.
// Excluded fields are id, created_at, deleted_at
func (d *Draft) Update() error {
	return d.ToContent().Update()
}

// Save creates or updates the content depending on the value of the id
func (d *Draft) Save() error {
	return d.ToContent().Save()
}

// FullyDelete removes an object from the database
func (d *Draft) FullyDelete() error {
	return d.ToContent().FullyDelete()
}

// Delete soft delete an object.
func (d *Draft) Delete() error {
	return d.ToContent().Delete()
}

// IsZero checks if the object is either nil or don't have an ID
func (d *Draft) IsZero() bool {
	return d.ToContent().IsZero()
}
