package router

// RouteAuth represents a middleware used to allow/block the access to an endpoint
type RouteAuth func(*Request) bool

// LoggedUser is a auth middleware that filters out anonymous users
func LoggedUser(req *Request) bool {
	return req.User != nil
}
