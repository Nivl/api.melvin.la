// +build integration

package api_test

import (
	"testing"

	"github.com/melvin-laplanche/ml-api/src/components/api"
	"github.com/stretchr/testify/assert"
)

func TestDefaultSetupCurrentEnv(t *testing.T) {
	args, deps, err := api.DefaultSetup()
	assert.NoError(t, err, "DefaultSetup() should have worked")
	assert.NotNil(t, args, "args should have been set")
	assert.NotNil(t, deps, "deps should have been set")
}
