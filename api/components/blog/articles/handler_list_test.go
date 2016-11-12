package articles_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Nivl/api.melvin.la/api/app/testhelpers"
	"github.com/Nivl/api.melvin.la/api/components/blog/articles"
	"github.com/stretchr/testify/assert"
)

// ListTest tests the List handler
func TestHandlerList(t *testing.T) {
	for i := 0; i < 10; i++ {
		a, u, s := articles.NewTestArticle(t, nil)
		testhelpers.SaveModel(t, a)
		testhelpers.SaveModel(t, u)
		testhelpers.SaveModel(t, s)
	}
	defer testhelpers.PurgeModels(t)

	tests := []struct {
		description string
		countWanted int
		code        int
	}{
		{"No params", 10, http.StatusOK},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			rec := callHandlerList(t)
			assert.Equal(t, tc.code, rec.Code)

			if rec.Code == http.StatusOK {
				var body *articles.PublicPayloads
				if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tc.countWanted, len(body.Results))
			}
		})
	}
}

func callHandlerList(t *testing.T) *httptest.ResponseRecorder {
	ri := &testhelpers.RequestInfo{
		Test:     t,
		Endpoint: articles.Endpoints[articles.EndpointList],
		URI:      "/blog/articles/",
	}

	return testhelpers.NewRequest(ri)
}
