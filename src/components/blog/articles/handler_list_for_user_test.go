package articles_test

// // TestHandlerListForUser tests the ListForUser handler
// func TestHandlerListForUser(t *testing.T) {
// 	defer testhelpers.PurgeModels(t)

// 	userNoArticles := authtest.NewUser(t, nil)
// 	userNoPrivate := authtest.NewUser(t, nil)
// 	userNoPublic := authtest.NewUser(t, nil)
// 	userAllKind := authtest.NewUser(t, nil)

// 	// Create 10 published articles
// 	for i := 0; i < 10; i++ {
// 		articlestest.NewArticle(t, &articles.Article{PublishedAt: db.Now(), UserID: userNoPrivate.ID})
// 		articlestest.NewArticle(t, &articles.Article{PublishedAt: db.Now(), UserID: userAllKind.ID})
// 	}

// 	// Create 10 unpublished articles
// 	for i := 0; i < 10; i++ {
// 		articlestest.NewArticle(t, &articles.Article{PublishedAt: nil, UserID: userNoPublic.ID})
// 		articlestest.NewArticle(t, &articles.Article{PublishedAt: nil, UserID: userAllKind.ID})
// 	}

// 	tests := []struct {
// 		description     string
// 		code            int
// 		authorID        string
// 		nbArticleWanted int
// 	}{
// 		{"Invalid author uuid", http.StatusBadRequest, "invalid uuid", 0},
// 		{"Un-existing author", http.StatusNotFound, "ab51713d-c2ce-4a10-afe8-38c264960a8e", 0},
// 		{"User without article", http.StatusOK, userNoArticles.ID, 0},
// 		{"User without public article", http.StatusOK, userNoPublic.ID, 0},
// 		{"User without private article", http.StatusOK, userNoPrivate.ID, 10},
// 		{"Regular user", http.StatusOK, userAllKind.ID, 10},
// 	}

// 	for i, tc := range tests {
// 		t.Run(tc.description, func(t *testing.T) {
// 			rec := callHandlerListForUser(t, tc.authorID)
// 			assert.Equal(t, tc.code, rec.Code, "Test %d", i)

// 			if testhelpers.Is2XX(rec.Code) {
// 				var body *articles.Payloads
// 				if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
// 					t.Fatal(err)
// 				}

// 				// Validate the number of articles
// 				assert.Equal(t, tc.nbArticleWanted, len(body.Results), "Test %d", i)

// 				// validate the content of the articles
// 				for _, art := range body.Results {
// 					// Validate the visibility of the article
// 					assert.NotNil(t, art.PublishedAt, "Test %d", i)

// 					// validate the author
// 					require.NotNil(t, art.User, "Test %d", i)
// 					assert.Equal(t, tc.authorID, art.User.ID, "Test %d", i)

// 					// Validate the content
// 					require.NotNil(t, art.Content, "Test %d", i)
// 					require.NotEmpty(t, art.Content.Title, "Test %d", i)
// 					assert.Nil(t, art.Draft, "Test %d", i)

// 				}
// 			}
// 		})
// 	}
// }

// func callHandlerListForUser(t *testing.T, authorID string) *httptest.ResponseRecorder {
// 	ri := &testhelpers.RequestInfo{
// 		Test:     t,
// 		Endpoint: articles.Endpoints[articles.EndpointUserList],
// 		URI:      fmt.Sprintf("/users/%s/blog/articles", authorID),
// 	}

// 	return testhelpers.NewRequest(ri)
// }
