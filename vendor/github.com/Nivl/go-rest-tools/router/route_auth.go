package router

// RouteAuth represents a middleware used to allow/block the access to an endpoint
type RouteAuth func(*Request) bool

// LoggedUserAccess is a auth middleware that filters out anonymous users
func LoggedUserAccess(req *Request) bool {
	return req.User != nil
}

// AdminAccess is a auth middleware that filters out non admin users
func AdminAccess(req *Request) bool {
	return req.User != nil && req.User.IsAdmin
}
