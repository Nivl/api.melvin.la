// +build integration

package education_test

import (
	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/melvin-laplanche/ml-api/src/components/api"
)

var (
	deps dependencies.Dependencies
)

func init() {
	var err error
	_, deps, err = api.DefaultSetup()
	if err != nil {
		panic(err)
	}
	httptests.DefaultRouter = api.GetRouter(deps)
}
