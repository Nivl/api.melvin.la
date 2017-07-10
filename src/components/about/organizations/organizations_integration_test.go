// +build integration

package organizations_test

import (
	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/melvin-laplanche/ml-api/src/components/api"
)

func init() {
	api.Setup()
	httptests.DefaultRouter = api.GetRouter()
}
