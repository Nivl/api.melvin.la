package router

// RouteHandler is the function signature we nee
type RouteHandler func(*Request, *Dependencies) error

// Endpoint represents an HTTP endpoint
type Endpoint struct {
	Verb string

	// Path is the path for the current component
	Path string

	// Auth is used to add a auth middleware
	Auth RouteAuth

	// Handler is the handler to call
	Handler RouteHandler

	// Params represents a list of params the endpoint needs
	Params interface{}
}
