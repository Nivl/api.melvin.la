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

	return New(p, pp)
}

// Paginator represents a pagination
type Paginator struct {
	currentPage int
	perPage     int
	MaxPerPage  int
}

// New creates a new Paginator
func New(page int, perPage int) *Paginator {
	return &Paginator{
		currentPage: page,
		perPage:     perPage,
		MaxPerPage:  100,
	}
}

// IsValid checks if the paginator is valid
func (p *Paginator) IsValid() bool {
	return (p.currentPage > 0) && (p.perPage <= p.MaxPerPage) && (p.perPage > 0)
}

// Offset returns a valid SQL offset value
func (p *Paginator) Offset() int {
	return (p.currentPage - 1) * p.perPage
}

// Limit returns a valid SQL limit value
func (p *Paginator) Limit() int {
	return p.perPage
}

// CurrentPage returns the current page index
func (p *Paginator) CurrentPage() int {
	return p.currentPage
}

// PerPage returns a the number of data a page should contain
func (p *Paginator) PerPage() int {
	return p.perPage
}
