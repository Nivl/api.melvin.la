package testguard

import (
	"net/url"
	"strings"
	"testing"

	"github.com/Nivl/go-rest-tools/router/guard"
	"github.com/stretchr/testify/assert"
)

// InvalidParamsTestCase represents the data needed to test params
type InvalidParamsTestCase struct {
	Description string
	MsgMatch    string
	Sources     map[string]url.Values
}

// InvalidParams checks if the params are correctly failing
func InvalidParams(t *testing.T, g *guard.Guard, testCases []InvalidParamsTestCase) {
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Description, func(t *testing.T) {
			t.Parallel()

			_, err := g.ParseParams(tc.Sources)
			if assert.Error(t, err, "expected the guard to fail") {
				assert.True(t, strings.Contains(err.Error(), tc.MsgMatch),
					"the error \"%s\" should contain the string \"%s\"", err.Error(), tc.MsgMatch)
			}
		})
	}
}
