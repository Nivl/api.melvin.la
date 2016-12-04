package articles

import "github.com/Nivl/sqalx"

// Draft Represents an article draft
type Draft Content

// ToContent returns a Content from a Draft
func (d *Draft) ToContent() *Content {
	content := Content(*d)
	return &content
}

// CreateTx persists a draft in the database using a transaction
func (d *Draft) CreateTx(tx sqalx.Node) error {
	return d.ToContent().CreateTx(tx)
}

// Create persists a draft in the database
func (d *Draft) Create() error {
	return d.ToContent().Create()
}

// UpdateTx updates most of the fields of a persisted content using a transaction
// Excluded fields are id, created_at, deleted_at
func (d *Draft) UpdateTx(tx sqalx.Node) error {
	return d.ToContent().UpdateTx(tx)
}

// Update updates most of the fields of a persisted content
// Excluded fields are id, created_at, deleted_at
func (d *Draft) Update() error {
	return d.ToContent().Update()
}

// SaveTx creates or updates the content depending on the value of the id using a transaction
func (d *Draft) SaveTx(tx sqalx.Node) error {
	return d.ToContent().SaveTx(tx)
}

// Save creates or updates the content depending on the value of the id
func (d *Draft) Save() error {
	return d.ToContent().Save()
}

// FullyDeleteTx removes an object from the database using a transaction
func (d *Draft) FullyDeleteTx(tx sqalx.Node) error {
	return d.ToContent().FullyDeleteTx(tx)
}

// FullyDelete removes an object from the database
func (d *Draft) FullyDelete() error {
	return d.ToContent().FullyDelete()
}

// DeleteTx soft delete an object using a transaction
func (d *Draft) DeleteTx(tx sqalx.Node) error {
	return d.ToContent().DeleteTx(tx)
}

// Delete soft delete an object
func (d *Draft) Delete() error {
	return d.ToContent().Delete()
}

// IsZero checks if the object is either nil or don't have an ID
func (d *Draft) IsZero() bool {
	return d.ToContent().IsZero()
}
