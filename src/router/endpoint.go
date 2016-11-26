package router

// RouteHandler is the function signature we nee
type RouteHandler func(*Request)

// Endpoint represents an HTTP endpoint
type Endpoint struct {
	Verb string

	// Path is the relative path for the current component
	// For `/v1/blog/articles/{id}` it would be `/articles/{id}``
	Path string

	// Prefix is used to add a prefix in front of the full path
	// Ex. the Prefix `/user/{uid}` can be use to make `/user/{uid}/blog/articles/{id}
	Prefix string

	// Auth is used to add a auth middleware
	Auth RouteAuth

	// Handler is the handler to call
	Handler RouteHandler

	// Params represents a list of params the endpoint needs
	Params interface{}
}
