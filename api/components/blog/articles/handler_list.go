package articles

import (
	"github.com/Nivl/api.melvin.la/api/router"
)

// HandlerList represents a API handler to get a list of articles
func HandlerList(req *router.Request) {
	arts := []*Article{}

	if err := Query().Find(defaultSearch).Sort("-createdAt").All(&arts); err != nil {
		req.Error(err)
		return
	}

	req.Ok(NewPayloadFromModels(arts))
}
