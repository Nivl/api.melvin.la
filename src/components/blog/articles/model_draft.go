package articles

// Draft Represents an article draft
type Draft Content

// Create persists a draft in the database
func (d *Draft) Create() error {
	content := Content(*d)
	return (&content).Create()
}

// Update updates most of the fields of a persisted content.
// Excluded fields are id, created_at, deleted_at
func (d *Draft) Update() error {
	content := Content(*d)
	return (&content).Update()
}

// Save creates or updates the content depending on the value of the id
func (d *Draft) Save() error {
	content := Content(*d)
	return (&content).Update()
}

// FullyDelete removes an object from the database
func (d *Draft) FullyDelete() error {
	content := Content(*d)
	return (&content).Update()
}

// Delete soft delete an object.
func (d *Draft) Delete() error {
	content := Content(*d)
	return (&content).Update()
}