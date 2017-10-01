// +build integration

package organizations_test

import (
	"os"
	"path"

	"github.com/Nivl/go-rest-tools/dependencies"
	"github.com/melvin-laplanche/ml-api/src/components/api"
)

var (
	migrationFolder string
)

func NewDeps() dependencies.Dependencies {
	var err error
	_, deps, err := api.DefaultSetup()
	if err != nil {
		panic(err)
	}
	return deps
}

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	migrationFolder = path.Join(wd, "..", "..", "..", "..", "db", "migrations")
}
