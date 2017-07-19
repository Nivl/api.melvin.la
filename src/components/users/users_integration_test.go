// +build integration

package users_test

import (
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/melvin-laplanche/ml-api/src/components/api"
)

var deps *router.Dependencies

func init() {
	api.Setup()
	deps, _ = router.NewDefaultDependencies()
	httptests.DefaultRouter = api.GetRouter()
}
