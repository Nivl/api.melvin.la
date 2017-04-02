package articles_test

// // TestHandlerUpdate tests the update handler
// func TestHandlerUpdate(t *testing.T) {
// 	defer testhelpers.PurgeModels(t)

// 	randomUser, randomUserSession := authtest.NewAuth(t)

// 	a1, u1, s1 := articlestest.NewArticle(t, nil)
// 	a2, u2, s2 := articlestest.NewArticle(t, nil)
// 	a3, u3, s3 := articlestest.NewArticle(t, nil)

// 	tests := []struct {
// 		description string
// 		code        int
// 		auth        *testhelpers.RequestAuth
// 		params      *articles.HandlerUpdateParams
// 		article     *articles.Article
// 	}{
// 		{
// 			"As anonymous",
// 			http.StatusUnauthorized,
// 			nil,
// 			&articles.HandlerUpdateParams{ID: a1.ID},
// 			a1,
// 		},
// 		{
// 			"As logged user",
// 			http.StatusNotFound,
// 			testhelpers.NewRequestAuth(randomUserSession.ID, randomUser.ID),
// 			&articles.HandlerUpdateParams{ID: a1.ID},
// 			a1,
// 		},
// 		{
// 			"Un-existing article",
// 			http.StatusNotFound,
// 			testhelpers.NewRequestAuth(randomUserSession.ID, randomUser.ID),
// 			&articles.HandlerUpdateParams{ID: "1228f726-54b0-4b28-ad8b-fa7d3c1c37b7"},
// 			nil,
// 		},
// 		{
// 			"Invalid article uuid",
// 			http.StatusBadRequest,
// 			testhelpers.NewRequestAuth(randomUserSession.ID, randomUser.ID),
// 			&articles.HandlerUpdateParams{ID: "invalid data"},
// 			nil,
// 		},
// 		{
// 			"As author - changing content",
// 			http.StatusOK,
// 			testhelpers.NewRequestAuth(s1.ID, u1.ID),
// 			&articles.HandlerUpdateParams{ID: a1.ID, Content: "New Content", Subtitle: "New Subtitle"},
// 			a1,
// 		},
// 		{
// 			"As author - wrong publish value",
// 			http.StatusBadRequest,
// 			testhelpers.NewRequestAuth(s2.ID, u2.ID),
// 			&articles.HandlerUpdateParams{ID: a2.ID, Publish: "maybe"},
// 			a2,
// 		},
// 		{
// 			"As author un-publish",
// 			http.StatusOK,
// 			testhelpers.NewRequestAuth(s2.ID, u2.ID),
// 			&articles.HandlerUpdateParams{ID: a2.ID, Publish: "false"},
// 			a2,
// 		},
// 		// needs to be below the un-publish as it's using the same article
// 		{
// 			"As author publish",
// 			http.StatusOK,
// 			testhelpers.NewRequestAuth(s2.ID, u2.ID),
// 			&articles.HandlerUpdateParams{ID: a2.ID, Publish: "true"},
// 			a2,
// 		},
// 		{
// 			"As author - invalid slug",
// 			http.StatusBadRequest,
// 			testhelpers.NewRequestAuth(s3.ID, u3.ID),
// 			&articles.HandlerUpdateParams{ID: a3.ID, Slug: "my slug"},
// 			a3,
// 		},
// 		{
// 			"As author - valid slug",
// 			http.StatusOK,
// 			testhelpers.NewRequestAuth(s3.ID, u3.ID),
// 			&articles.HandlerUpdateParams{ID: a3.ID, Slug: "my-slug"},
// 			a3,
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.description, func(t *testing.T) {
// 			rec := callHandlerUpdate(t, tc.params, tc.auth)
// 			assert.Equal(t, tc.code, rec.Code)

// 			if testhelpers.Is2XX(rec.Code) {
// 				var pld *articles.Payload
// 				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
// 					t.Fatal(err)
// 				}

// 				if tc.params.Slug != "" {
// 					assert.Equal(t, tc.params.Slug, pld.Slug)
// 				} else {
// 					assert.Equal(t, tc.article.Slug, pld.Slug)
// 				}

// 				require.NotNil(t, pld.Content)

// 				if tc.params.Title != "" {
// 					assert.Equal(t, tc.params.Title, pld.Content.Title)
// 				} else {
// 					assert.Equal(t, tc.article.Content.Title, pld.Content.Title)
// 				}

// 				if tc.params.Subtitle != "" {
// 					assert.Equal(t, tc.params.Subtitle, pld.Content.Subtitle)
// 				} else {
// 					assert.Equal(t, tc.article.Content.Subtitle, pld.Content.Subtitle)
// 				}

// 				if tc.params.Description != "" {
// 					assert.Equal(t, tc.params.Description, pld.Content.Description)
// 				} else {
// 					assert.Equal(t, tc.article.Content.Description, pld.Content.Description)
// 				}

// 				if tc.params.Content != "" {
// 					assert.Equal(t, tc.params.Content, pld.Content.Content)
// 				} else {
// 					assert.Equal(t, tc.article.Content.Content, pld.Content.Content)
// 				}

// 				if tc.params.Publish == "true" {
// 					assert.NotNil(t, pld.PublishedAt)
// 				} else if tc.params.Publish == "false" {
// 					assert.Nil(t, pld.PublishedAt)
// 				}
// 			}
// 		})
// 	}
// }

// func callHandlerUpdate(t *testing.T, params *articles.HandlerUpdateParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
// 	ri := &testhelpers.RequestInfo{
// 		Test:     t,
// 		Endpoint: articles.Endpoints[articles.EndpointUpdate],
// 		URI:      fmt.Sprintf("/blog/articles/%s", params.ID),
// 		Auth:     auth,
// 		Params:   params,
// 	}

// 	return testhelpers.NewRequest(ri)
// }
