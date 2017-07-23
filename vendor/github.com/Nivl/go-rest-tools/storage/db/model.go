package db

// Model represents a database model
type Model interface {
	// GetID returns the ID of the model. We use the "Get" word to avoid
	// collision with the ID field of the models (and to stay backward compatible)
	GetID() string
	// SetID sets the ID of the model
	SetID(string)
}
