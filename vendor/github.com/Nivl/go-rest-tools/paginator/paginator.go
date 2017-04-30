package paginator

// HandlerParams represents the params needed to handle a pagination
type HandlerParams struct {
	// Page represents the number of the current page
	Page *int `from:"query" json:"page" default:"1"`

	// PerPage represents the maximum number of element we want per page
	// A value of -1 will
	PerPage *int `from:"query" json:"per_page"`
}

// Paginator returns a Paginator from an HandlerParams
func (params *HandlerParams) Paginator(defaultPerPage int) *Paginator {
	// Should never be nil because of the default value
	p := *params.Page

	pp := defaultPerPage
	if params.PerPage != nil {
		pp = *params.PerPage
	}

	return &Paginator{
		CurrentPage: p,
		PerPage:     pp,
		MaxPerPage:  100,
	}
}

// Paginator represents a pagination
type Paginator struct {
	CurrentPage int
	PerPage     int
	MaxPerPage  int
}

// New creates a new Paginator
func New(page int, perPage int) *Paginator {
	return &Paginator{
		CurrentPage: page,
		PerPage:     perPage,
		MaxPerPage:  100,
	}
}

// IsValid checks if the paginator is valid
func (p *Paginator) IsValid() bool {
	return (p.CurrentPage > 0) && (p.PerPage <= p.MaxPerPage) && (p.PerPage > 0)
}

// Offset returns a valid SQL offset value
func (p *Paginator) Offset() int {
	return (p.CurrentPage - 1) * p.PerPage
}

// Limit returns a valid SQL limit value
func (p *Paginator) Limit() int {
	return p.PerPage
}
