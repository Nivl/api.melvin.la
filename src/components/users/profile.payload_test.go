package users_test

import (
	"testing"

	"github.com/melvin-laplanche/ml-api/src/components/users"
	"github.com/stretchr/testify/assert"
)

func TestExportPublicProfileNil(t *testing.T) {
	var exp *users.Profile
	pld := exp.ExportPublic()
	assert.Nil(t, pld, "nil should export as nil")
}

func TestExportPrivateProfileNil(t *testing.T) {
	var exp *users.Profile
	pld := exp.ExportPrivate()
	assert.Nil(t, pld, "nil should export as nil")
}
