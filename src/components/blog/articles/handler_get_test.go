package articles_test

// func TestHandlerGet(t *testing.T) {
// 	defer testhelpers.PurgeModels(t)

// 	// Random logged user
// 	randomUser, randomUserSession := authtest.NewAuth(t)

// 	// Unpublished article with no draft
// 	a, uA, sA := articlestest.NewArticle(t, &articles.Article{PublishedAt: nil})

// 	// Published article with no draft
// 	aPublished, uPublished, sPublished := articlestest.NewArticle(t, nil)

// 	// Published article with draft
// 	aPubDraft, uPubDraft, sPubDraft := articlestest.NewArticle(t,
// 		&articles.Article{PublishedAt: db.Now(),
// 			Content: &articles.Content{Title: "Title", Content: "Content"},
// 			Draft:   &articles.Draft{Title: "Title Draft", Content: "Content Draft"},
// 		})

// 	tests := []struct {
// 		description string
// 		code        int
// 		params      *articles.HandlerGetParams
// 		auth        *testhelpers.RequestAuth
// 		a           *articles.Article
// 		hasDraft    bool
// 	}{
// 		// not published - no draft
// 		{
// 			"Anonymous - not published - no draft",
// 			http.StatusNotFound,
// 			&articles.HandlerGetParams{ID: a.ID},
// 			nil,
// 			a, false,
// 		},
// 		{
// 			"logged user - not published - no draft",
// 			http.StatusNotFound,
// 			&articles.HandlerGetParams{ID: a.ID},
// 			testhelpers.NewRequestAuth(randomUserSession.ID, randomUser.ID),
// 			a, false,
// 		},
// 		{
// 			"Author - not published - no draft",
// 			http.StatusOK,
// 			&articles.HandlerGetParams{ID: a.ID},
// 			testhelpers.NewRequestAuth(sA.ID, uA.ID),
// 			a, false,
// 		},

// 		// published - no draft
// 		{
// 			"Anonymous - published - no draft",
// 			http.StatusOK,
// 			&articles.HandlerGetParams{ID: aPublished.ID},
// 			nil,
// 			aPublished, false,
// 		},
// 		{
// 			"logged user - published - no draft",
// 			http.StatusOK,
// 			&articles.HandlerGetParams{ID: aPublished.ID},
// 			testhelpers.NewRequestAuth(randomUserSession.ID, randomUser.ID),
// 			aPublished, false,
// 		},
// 		{
// 			"Author - published - no draft",
// 			http.StatusOK,
// 			&articles.HandlerGetParams{ID: aPublished.ID},
// 			testhelpers.NewRequestAuth(sPublished.ID, uPublished.ID),
// 			aPublished, false,
// 		},

// 		// published - draft
// 		{
// 			"Anonymous - published - draft",
// 			http.StatusOK,
// 			&articles.HandlerGetParams{ID: aPubDraft.ID},
// 			nil,
// 			aPubDraft, false,
// 		},
// 		{
// 			"logged user - published - draft",
// 			http.StatusOK,
// 			&articles.HandlerGetParams{ID: aPubDraft.ID},
// 			testhelpers.NewRequestAuth(randomUserSession.ID, randomUser.ID),
// 			aPubDraft, false,
// 		},
// 		{
// 			"Author - published - draft",
// 			http.StatusOK,
// 			&articles.HandlerGetParams{ID: aPubDraft.ID},
// 			testhelpers.NewRequestAuth(sPubDraft.ID, uPubDraft.ID),
// 			aPubDraft, true,
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.description, func(t *testing.T) {
// 			rec := callHandlerGet(t, tc.params, tc.auth)
// 			assert.Equal(t, tc.code, rec.Code)

// 			if testhelpers.Is2XX(rec.Code) {
// 				var pld articles.Payload
// 				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
// 					t.Fatal(err)
// 				}

// 				assert.Equal(t, tc.a.ID, pld.ID)
// 				assert.Equal(t, tc.a.Slug, pld.Slug)

// 				require.NotNil(t, pld.User)
// 				assert.Equal(t, tc.a.UserID, pld.User.ID)

// 				assert.Equal(t, tc.a.Content.Title, pld.Content.Title)
// 				assert.Equal(t, tc.a.Content.Subtitle, pld.Content.Subtitle)
// 				assert.Equal(t, tc.a.Content.Description, pld.Content.Description)
// 				assert.Equal(t, tc.a.Content.Content, pld.Content.Content)

// 				if tc.hasDraft {
// 					require.NotNil(t, pld.Draft)
// 					assert.Equal(t, tc.a.Draft.Title, pld.Draft.Title)
// 					assert.Equal(t, tc.a.Draft.Subtitle, pld.Draft.Subtitle)
// 					assert.Equal(t, tc.a.Draft.Description, pld.Draft.Description)
// 					assert.Equal(t, tc.a.Draft.Content, pld.Draft.Content)
// 				} else {
// 					assert.Nil(t, pld.Draft)
// 				}
// 			}
// 		})
// 	}
// }

// func callHandlerGet(t *testing.T, params *articles.HandlerGetParams, auth *testhelpers.RequestAuth) *httptest.ResponseRecorder {
// 	ri := &testhelpers.RequestInfo{
// 		Test:     t,
// 		Endpoint: articles.Endpoints[articles.EndpointGet],
// 		URI:      fmt.Sprintf("/blog/articles/%s", params.ID),
// 		Params:   params,
// 		Auth:     auth,
// 	}

// 	return testhelpers.NewRequest(ri)
// }
