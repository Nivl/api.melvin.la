package articles_test

// func TestHandlerDelete(t *testing.T) {
// 	defer testhelpers.PurgeModels(t)

// 	// Random logged user
// 	randomUser, randomUserSession := authtest.NewAuth(t)

// 	// Published article with no draft
// 	a, u, s := articlestest.NewArticle(t, nil)

// 	tests := []struct {
// 		description string
// 		code        int
// 		params      *articles.HandlerDeleteParams
// 		auth        *testhelpers.RequestAuth
// 	}{
// 		{
// 			"Anonymous",
// 			http.StatusUnauthorized,
// 			&articles.HandlerDeleteParams{ID: a.ID},
// 			nil,
// 		},
// 		{
// 			"logged user",
// 			http.StatusNotFound,
// 			&articles.HandlerDeleteParams{ID: a.ID},
// 			testhelpers.NewRequestAuth(randomUserSession.ID, randomUser.ID),
// 		},
// 		{
// 			"Owner",
// 			http.StatusNoContent,
// 			&articles.HandlerDeleteParams{ID: a.ID},
// 			testhelpers.NewRequestAuth(s.ID, u.ID),
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.description, func(t *testing.T) {
// 			rec := callHandlerDelete(t, tc.params, tc.auth)
// 			assert.Equal(t, tc.code, rec.Code)

// 			if testhelpers.Is2XX(rec.Code) {
// 				exists, err := articles.Exists(tc.params.ID)
// 				if err != nil {
// 					t.Fatal(err)
// 				}
// 				assert.False(t, exists)
// 			}
// 		})
// 	}
// }

// func callHandlerDelete(t *testing.T, params *articles.HandlerDeleteParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
// 	ri := &testhelpers.RequestInfo{
// 		Test:     t,
// 		Endpoint: articles.Endpoints[articles.EndpointDelete],
// 		URI:      fmt.Sprintf("/blog/articles/%s", params.ID),
// 		Params:   params,
// 		Auth:     auth,
// 	}

// 	return testhelpers.NewRequest(ri)
// }
