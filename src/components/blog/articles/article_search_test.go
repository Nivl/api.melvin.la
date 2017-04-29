package articles_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Nivl/go-rest-tools/network/http/httptests"
	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/primitives/models/lifecycle"
	"github.com/Nivl/go-rest-tools/primitives/ptrs"
	"github.com/Nivl/go-rest-tools/security/auth/testdata"
	"github.com/Nivl/go-rest-tools/storage/db"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles"
	"github.com/melvin-laplanche/ml-api/src/components/blog/articles/articletestdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	testCases := []struct {
		description string
		f           func(t *testing.T)
	}{
		{"wrong input", searchWrongInput},
		{"ordering", searchOrdering},
		{"pagination", searchPagination},
		{"status", searchStatus},
		{"full text", searchFullText},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.description, tc.f)
	}
}

func searchWrongInput(t *testing.T) {
	t.Parallel()
	defer lifecycle.PurgeModels(t)

	t.Run("parallel wrapper", func(t *testing.T) {
		testCases := []struct {
			description string
			code        int
			params      *articles.SearchParams
		}{
			{
				"Access page -1",
				http.StatusBadRequest,
				&articles.SearchParams{
					HandlerParams: paginator.HandlerParams{Page: ptrs.NewInt(-1)},
				},
			},
			{
				"Get 10000000 article",
				http.StatusBadRequest,
				&articles.SearchParams{
					HandlerParams: paginator.HandlerParams{PerPage: ptrs.NewInt(100000)},
				},
			},
			{
				"Get 0 articles",
				http.StatusBadRequest,
				&articles.SearchParams{
					HandlerParams: paginator.HandlerParams{PerPage: ptrs.NewInt(0)},
				},
			},
			{
				"Order by a",
				http.StatusBadRequest,
				&articles.SearchParams{OrderBy: "a"},
			},
			{
				"Order by non allowed field",
				http.StatusBadRequest,
				&articles.SearchParams{OrderBy: "deleted_at"},
			},
			{
				"Filter by non allowed status",
				http.StatusBadRequest,
				&articles.SearchParams{OrderBy: "deleted_at"},
			},
			{
				"Filter by non allowed status",
				http.StatusBadRequest,
				&articles.SearchParams{Status: "deleted_at"},
			},
			{
				"Filter admin only status",
				http.StatusBadRequest,
				&articles.SearchParams{Status: "deleted"},
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()

				rec := callSearch(t, tc.params, nil)
				require.Equal(t, tc.code, rec.Code)
			})
		}
	})
}

func searchPagination(t *testing.T) {
	defer lifecycle.PurgeModels(t)

	arts := []*articles.Article{}
	for i := 0; i < articles.DefaultNbResultsPerPage*2.5; i++ {
		futureTime := time.Now().Add(time.Duration(i+i) * -time.Second)
		future := &db.Time{Time: futureTime}
		base := &articles.Article{CreatedAt: future, PublishedAt: future}
		a, _, _ := articletestdata.NewArticle(t, base)
		arts = append(arts, a)
	}

	t.Run("parallel wrapper", func(t *testing.T) {
		testCases := []struct {
			description string
			params      *articles.SearchParams
			total       int
		}{
			{
				"Page 1 full",
				&articles.SearchParams{
					HandlerParams: paginator.HandlerParams{Page: ptrs.NewInt(1)},
				},
				articles.DefaultNbResultsPerPage,
			},
			{
				"Page 2 full",
				&articles.SearchParams{
					HandlerParams: paginator.HandlerParams{Page: ptrs.NewInt(2)},
				},
				articles.DefaultNbResultsPerPage,
			},
			{
				"Page 3 half",
				&articles.SearchParams{
					HandlerParams: paginator.HandlerParams{Page: ptrs.NewInt(3)},
				},
				articles.DefaultNbResultsPerPage / 2,
			},
			{
				"Page 1 limit 5",
				&articles.SearchParams{
					HandlerParams: paginator.HandlerParams{Page: ptrs.NewInt(1), PerPage: ptrs.NewInt(5)},
				},
				5,
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callSearch(t, tc.params, nil)
				require.Equal(t, http.StatusOK, rec.Code)

				var pld *articles.Payloads
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.total, len(pld.Results))

				perPage := articles.DefaultNbResultsPerPage
				if tc.params.PerPage != nil {
					perPage = *tc.params.Page
				}
				delta := (*tc.params.Page - 1) * perPage

				for i, a := range pld.Results {
					assert.Equal(t, arts[i+delta].ID, a.ID)
				}
			})
		}
	})
}

func searchOrdering(t *testing.T) {
	defer lifecycle.PurgeModels(t)

	asc := []*articles.Article{}
	desc := []*articles.Article{}
	nbArticles := articles.DefaultNbResultsPerPage
	for i := 0; i < nbArticles; i++ {
		futureTime := time.Now().Add(time.Duration(i+i) * time.Second)
		future := &db.Time{Time: futureTime}
		base := &articles.Article{CreatedAt: future, PublishedAt: future}
		a, _, _ := articletestdata.NewArticle(t, base)
		asc = append(asc, a)
	}

	for i := nbArticles - 1; i > -1; i-- {
		desc = append(desc, asc[i])
	}

	t.Run("parallel wrapper", func(t *testing.T) {
		testCases := []struct {
			description string
			params      *articles.SearchParams
			resWanted   []*articles.Article
		}{
			{
				"published_at Descending",
				&articles.SearchParams{
					HandlerParams: paginator.HandlerParams{Page: ptrs.NewInt(1)},
					OrderBy:       "-published_at",
				},
				desc,
			},
			{
				"published_at Ascending",
				&articles.SearchParams{
					HandlerParams: paginator.HandlerParams{Page: ptrs.NewInt(1)},
					OrderBy:       "published_at",
				},
				asc,
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callSearch(t, tc.params, nil)
				require.Equal(t, http.StatusOK, rec.Code)

				var pld *articles.Payloads
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}

				require.Equal(t, len(tc.resWanted), len(pld.Results))

				for i, a := range pld.Results {
					assert.Equal(t, tc.resWanted[i].ID, a.ID)
				}
			})
		}
	})
}

func searchStatus(t *testing.T) {
	defer lifecycle.PurgeModels(t)

	nbPublished := 5
	for i := 0; i < nbPublished; i++ {
		futureTime := time.Now().Add(time.Duration(i+i) * time.Second)
		future := &db.Time{Time: futureTime}
		base := &articles.Article{CreatedAt: future, PublishedAt: future}
		articletestdata.NewArticle(t, base)
	}

	nbUnPublished := 4
	for i := 0; i < nbUnPublished; i++ {
		base := &articles.Article{PublishedAt: nil}
		articletestdata.NewArticle(t, base)
	}

	nbDeleted := 3
	for i := 0; i < nbDeleted; i++ {
		base := &articles.Article{PublishedAt: nil, DeletedAt: db.Now()}
		articletestdata.NewArticle(t, base)
	}

	_, adminSession := testdata.NewAdminAuth(t)
	adminAuth := httptests.NewRequestAuth(adminSession)

	t.Run("parallel wrapper", func(t *testing.T) {
		testCases := []struct {
			description string
			params      *articles.SearchParams
			totalWanted int
		}{
			{
				"Published article",
				&articles.SearchParams{Status: "published"},
				nbPublished,
			},
			{
				"UnPublished article",
				&articles.SearchParams{Status: "unpublished"},
				nbUnPublished,
			},
			{
				"Deleted article",
				&articles.SearchParams{Status: "deleted"},
				nbDeleted,
			},
			{
				"Published AND unPublished article",
				&articles.SearchParams{Status: "published|unpublished"},
				nbPublished + nbUnPublished,
			},
			{
				"unPublished AND deleted article",
				&articles.SearchParams{Status: "unpublished|deleted"},
				nbUnPublished + nbDeleted,
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callSearch(t, tc.params, adminAuth)
				require.Equal(t, http.StatusOK, rec.Code)

				var pld *articles.Payloads
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.totalWanted, len(pld.Results))
			})
		}
	})
}

func searchFullText(t *testing.T) {
	defer lifecycle.PurgeModels(t)

	// Articles with Californication
	californication := []*articles.Article{
		&articles.Article{PublishedAt: db.Now(), Version: &articles.Version{
			Content: "Californication Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque feugiat nisi a felis fermentum dapibus. Fusce eleifend arcu nec leo posuere, vel vestibulum nibh tincidunt. Curabitur sagittis sagittis dolor id dapibus. Aliquam erat volutpat. Maecenas in accumsan ligula. Morbi augue justo, dictum quis nisi sed, malesuada ornare neque. Praesent id vehicula sem. Integer in urna sagittis, pretium massa fringilla, vestibulum eros. Vivamus ullamcorper diam non dictum tincidunt. Nunc et pellentesque nibh, vitae rutrum dolor. Vivamus quis semper lectus. In tincidunt lorem ex. Mauris pharetra pharetra mauris eu luctus. Etiam mauris velit, ultrices vitae malesuada nec, tincidunt aliquet purus. Mauris gravida dui a ullamcorper commodo. Vivamus eget dolor quis ipsum pellentesque bibendum.",
		}},
		&articles.Article{PublishedAt: db.Now(), Version: &articles.Version{
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque Californication  feugiat nisi a felis fermentum dapibus. Fusce eleifend arcu nec leo posuere, vel vestibulum nibh tincidunt. Curabitur sagittis sagittis dolor id dapibus. Aliquam erat volutpat. Maecenas in accumsan ligula. Morbi augue justo, dictum quis nisi sed, malesuada ornare neque. Praesent id vehicula sem. Integer in urna sagittis, pretium massa fringilla, vestibulum eros. Vivamus ullamcorper diam non dictum tincidunt. Nunc et pellentesque nibh, vitae rutrum dolor. Vivamus quis semper lectus. In tincidunt lorem ex. Mauris pharetra pharetra mauris eu luctus. Etiam mauris velit, ultrices vitae malesuada nec, tincidunt aliquet purus. Mauris gravida dui a ullamcorper commodo. Vivamus eget dolor quis ipsum pellentesque bibendum.",
		}},
		&articles.Article{PublishedAt: db.Now(), Version: &articles.Version{
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque  feugiat nisi a felis fermentum dapibus. Fusce eleifend arcu nec leo posuere, vel vestibulum nibh tincidunt. Curabitur sagittis sagittis dolor id dapibus. Aliquam erat volutpat. Maecenas in accumsan ligula. Morbi augue justo, dictum quis nisi sed, malesuada ornare neque. Praesent id vehicula sem. Integer in urna sagittis, pretium massa fringilla, vestibulum eros. Vivamus ullamcorper diam non dictum tincidunt. Nunc et pellentesque nibh, vitae rutrum dolor. Vivamus quis semper lectus. In tincidunt lorem ex. Mauris pharetra pharetra mauris eu luctus. Etiam mauris velit, ultrices vitae malesuada nec, tincidunt aliquet purus. Mauris gravida dui a ullamcorper commodo. Vivamus eget dolor quis ipsum pellentesque bibendum. Californication",
		}},
		&articles.Article{PublishedAt: db.Now(), Version: &articles.Version{
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque  feugiat nisi a felis fermentum dapibus. Fusce eleifend arcu nec leo posuere, vel vestibulum nibh tincidunt. Curabitur sagittis sagittis dolor id dapibus. Aliquam erat volutpat. Maecenas in accumsan ligula. Morbi augue justo, dictum quis nisi sed, malesuada ornare neque. Praesent id vehicula sem. Integer in urna sagittis, pretium massa fringilla, vestibulum eros. Vivamus ullamcorper diam non dictum tincidunt. Nunc et pellentesque nibh, vitae rutrum dolor. Vivamus quis semper lectus. In tincidunt lorem ex. Mauris pharetra pharetra mauris eu luctus. Etiam mauris velit, ultrices vitae malesuada nec, tincidunt aliquet purus. Mauris gravida dui a ullamcorper Californication commodo. Vivamus eget dolor quis ipsum pellentesque bibendum.",
		}},
		&articles.Article{PublishedAt: db.Now(), Version: &articles.Version{
			Title: "Californication",
		}},
		&articles.Article{PublishedAt: db.Now(), Version: &articles.Version{
			Title: "Lorem ipsum dolor Californication sit amet",
		}},
	}
	for _, a := range californication {
		articletestdata.NewArticle(t, a)
	}

	// Articles with Bates Motel
	batesMotel := []*articles.Article{
		&articles.Article{PublishedAt: db.Now(), Version: &articles.Version{
			Content: "Bates Motel Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque feugiat nisi a felis fermentum dapibus. Fusce eleifend arcu nec leo posuere, vel vestibulum nibh tincidunt. Curabitur sagittis sagittis dolor id dapibus. Aliquam erat volutpat. Maecenas in accumsan ligula. Morbi augue justo, dictum quis nisi sed, malesuada ornare neque. Praesent id vehicula sem. Integer in urna sagittis, pretium massa fringilla, vestibulum eros. Vivamus ullamcorper diam non dictum tincidunt. Nunc et pellentesque nibh, vitae rutrum dolor. Vivamus quis semper lectus. In tincidunt lorem ex. Mauris pharetra pharetra mauris eu luctus. Etiam mauris velit, ultrices vitae malesuada nec, tincidunt aliquet purus. Mauris gravida dui a ullamcorper commodo. Vivamus eget dolor quis ipsum pellentesque bibendum.",
		}},
		&articles.Article{PublishedAt: db.Now(), Version: &articles.Version{
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque Bates Motel  feugiat nisi a felis fermentum dapibus. Fusce eleifend arcu nec leo posuere, vel vestibulum nibh tincidunt. Curabitur sagittis sagittis dolor id dapibus. Aliquam erat volutpat. Maecenas in accumsan ligula. Morbi augue justo, dictum quis nisi sed, malesuada ornare neque. Praesent id vehicula sem. Integer in urna sagittis, pretium massa fringilla, vestibulum eros. Vivamus ullamcorper diam non dictum tincidunt. Nunc et pellentesque nibh, vitae rutrum dolor. Vivamus quis semper lectus. In tincidunt lorem ex. Mauris pharetra pharetra mauris eu luctus. Etiam mauris velit, ultrices vitae malesuada nec, tincidunt aliquet purus. Mauris gravida dui a ullamcorper commodo. Vivamus eget dolor quis ipsum pellentesque bibendum.",
		}},
		&articles.Article{PublishedAt: db.Now(), Version: &articles.Version{
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque  feugiat nisi a felis fermentum dapibus. Fusce eleifend arcu nec leo posuere, vel vestibulum nibh tincidunt. Curabitur sagittis sagittis dolor id dapibus. Aliquam erat volutpat. Maecenas in accumsan ligula. Morbi augue justo, dictum quis nisi sed, malesuada ornare neque. Praesent id vehicula sem. Integer in urna sagittis, pretium massa fringilla, vestibulum eros. Vivamus ullamcorper diam non dictum tincidunt. Nunc et pellentesque nibh, vitae rutrum dolor. Vivamus quis semper lectus. In tincidunt lorem ex. Mauris pharetra pharetra mauris eu luctus. Etiam mauris velit, ultrices vitae malesuada nec, tincidunt aliquet purus. Mauris gravida dui a ullamcorper commodo. Vivamus eget dolor quis ipsum pellentesque bibendum. Bates Motel",
		}},
		&articles.Article{PublishedAt: db.Now(), Version: &articles.Version{
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque  feugiat nisi a felis fermentum dapibus. Fusce eleifend arcu nec leo posuere, vel vestibulum nibh tincidunt. Curabitur sagittis sagittis dolor id dapibus. Aliquam erat volutpat. Maecenas in accumsan ligula. Morbi augue justo, dictum quis nisi sed, malesuada ornare neque. Praesent id vehicula sem. Integer in urna sagittis, pretium massa fringilla, vestibulum eros. Vivamus ullamcorper diam non dictum tincidunt. Nunc et pellentesque nibh, vitae rutrum dolor. Vivamus quis semper lectus. In tincidunt lorem ex. Mauris pharetra pharetra mauris eu luctus. Etiam mauris velit, ultrices vitae malesuada nec, tincidunt aliquet purus. Mauris gravida dui a ullamcorper Bates Motel commodo. Vivamus eget dolor quis ipsum pellentesque bibendum.",
		}},
		&articles.Article{PublishedAt: db.Now(), Version: &articles.Version{
			Content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Quisque  feugiat nisi a felis fermentum dapibus. Fusce eleifend arcu nec leo posuere, vel vestibulum nibh tincidunt. Curabitur sagittis sagittis dolor id dapibus. Aliquam erat volutpat. Maecenas in accumsan ligula. Morbi augue justo, dictum quis nisi sed, malesuada Bates Motel ornare neque. Praesent id vehicula sem. Integer in urna sagittis, pretium massa fringilla, vestibulum eros. Vivamus ullamcorper diam non dictum tincidunt. Nunc et pellentesque nibh, vitae rutrum dolor. Vivamus quis semper lectus. In tincidunt lorem ex. Mauris pharetra pharetra mauris eu luctus. Etiam mauris velit, ultrices vitae malesuada nec, tincidunt aliquet purus. Mauris gravida dui a ullamcorper commodo. Vivamus eget dolor quis ipsum pellentesque bibendum.",
		}},
	}
	for _, a := range batesMotel {
		articletestdata.NewArticle(t, a)
	}

	t.Run("parallel wrapper", func(t *testing.T) {
		testCases := []struct {
			description string
			params      *articles.SearchParams
			totalWanted int
		}{
			{
				"articles with Californication",
				&articles.SearchParams{Query: "californication"},
				len(californication),
			},
			{
				"articles with Bates Motel",
				&articles.SearchParams{Query: "Bates Motel"},
				len(batesMotel),
			},
		}

		for _, tc := range testCases {
			tc := tc
			t.Run(tc.description, func(t *testing.T) {
				t.Parallel()
				rec := callSearch(t, tc.params, nil)
				require.Equal(t, http.StatusOK, rec.Code)

				var pld *articles.Payloads
				if err := json.NewDecoder(rec.Body).Decode(&pld); err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.totalWanted, len(pld.Results))
			})
		}
	})
}

func callSearch(t *testing.T, params *articles.SearchParams, auth *httptests.RequestAuth) *httptest.ResponseRecorder {
	ri := &httptests.RequestInfo{
		Endpoint: articles.Endpoints[articles.EndpointSearch],
		Auth:     auth,
		Params:   params,
	}

	return httptests.NewRequest(t, ri)
}
