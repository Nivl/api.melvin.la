package router

type RouteAuth func(*Request) bool

func LoggedUser(req *Request) bool {
	return req.User != nil
}
