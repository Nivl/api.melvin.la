package articles

import (
	"fmt"
	"strings"

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/paginator"
	"github.com/Nivl/go-rest-tools/primitives/slices"
	"github.com/Nivl/go-rest-tools/router"
	"github.com/Nivl/go-rest-tools/security/auth"
	"github.com/Nivl/go-rest-tools/storage/db"
)

// DefaultNbResultsPerPage represents the default number of result per page
const DefaultNbResultsPerPage = 20

// SearchParams represents the params needed by the Search handler
type SearchParams struct {
	paginator.HandlerParams

	// Query represents a string to use to look against the content field
	Query string `from:"query" json:"q" params:"trim"`

	// OrderBy represents a list of orders separated by "|"
	// ex ?order=-name|day_of_week will order by name desc and day of week asc
	// Accepted fields are: created_at, published_at
	OrderBy string `from:"query" json:"order" params:"trim"`

	// Status represents a list of article status separated by "|"
	// ex ?status=published|deleted will return both published and deleted
	// articles
	// Accepted fields are: published, deleted, unpublished
	Status string `from:"query" json:"status" params:"trim" default:"published"`
}

// Search represents an API handler to get a list of articles
func Search(req *router.Request) error {
	params := req.Params.(*SearchParams)

	pagination := params.Paginator(DefaultNbResultsPerPage)
	if !pagination.IsValid() {
		return httperr.NewBadRequest("Invalid pagination data")
	}

	// Set default SQL params
	selct := "articles.*"
	join := ""
	where := ""
	orderBy := ""
	args := map[string]interface{}{}

	// Full text search on the content
	if params.Query != "" {
		selct += ", ts_rank(version.search_data_vector, keywords, 1) AS rank"
		join += "CROSS JOIN plainto_tsquery(:search_query) keywords"
		where = sqlSetComma(where, "keywords @@ version.search_data_vector")
		orderBy = "rank DESC"
		args["search_query"] = params.Query
	}

	order, err := params.ParseOrder()
	if err != nil {
		return err
	}
	orderBy = sqlSetComma(orderBy, order)

	status, err := params.ParseStatus(req.User)
	if err != nil {
		return err
	}
	where = sqlSetAnd(where, status)

	// Set SQL keywords if needed
	if where != "" {
		where = " WHERE " + where
	}
	// We default the ordering by recency
	if orderBy == "" {
		orderBy = "published_at DESC, created_at DESC"
	}

	// Exec query and return payload
	var list Articles
	stmt := `SELECT %s,
						` + auth.UserJoinSQL("users") + `,
						` + JoinVersionSQL("version") + `
						FROM blog_articles articles
						JOIN users ON users.id = articles.user_id
						JOIN blog_article_versions version ON version.id = articles.current_version
						%s
						%s
						ORDER BY %s LIMIT %d OFFSET %d`
	stmt = fmt.Sprintf(stmt, selct, join, where, orderBy, pagination.Limit(), pagination.Offset())
	if err := db.NamedSelect(&list, stmt, args); err != nil {
		return err
	}

	req.Ok(list.Export())
	return nil
}

// ParseOrder parse the OrderBy fields and returns the associated ORDER BY clause
func (params *SearchParams) ParseOrder() (string, error) {
	orderBy := ""

	if params.OrderBy != "" {
		sortableFields := []string{"created_at", "published_at"}
		fields := strings.Split(params.OrderBy, "|")
		for _, f := range fields {
			f = strings.ToLower(f)
			// we need at least 2 chars (ex. -a)
			if len(f) < 2 {
				return "", httperr.NewBadRequest("invalid sort option: %s", f)
			}
			order := "ASC"
			if f[0] == '-' {
				order = "DESC"
				f = f[1:]
			}
			found, err := slices.InSlice(sortableFields, f)
			if err != nil {
				return "", err
			}
			if !found {
				return "", httperr.NewBadRequest("field not sortable: %s", f)
			}
			if orderBy != "" {
				orderBy += ", "
			}
			orderBy += fmt.Sprintf("%s %s", f, order)
		}
	}
	return orderBy, nil
}

// ParseStatus parse the status fields and returns the associated WHERE clause
func (params *SearchParams) ParseStatus(u *auth.User) (string, error) {
	where := ""
	var err error
	fields := strings.Split(params.Status, "|")
	statuses := struct {
		published   bool
		unpublished bool
		deleted     bool
	}{}

	// we validate the fields
	optsAllowed := []string{"published", "unpublished", "deleted"}
	for _, f := range fields {
		f = strings.ToLower(f)
		found, err := slices.InSlice(optsAllowed, f)
		if err != nil {
			return "", err
		}
		if !found {
			return "", httperr.NewBadRequest("field not filterable: %s", f)
		}
	}

	// public statuses
	statuses.published, err = slices.InSlice(fields, "published")
	if err != nil {
		return "", err
	}
	// private statuses, we only set them for admin users
	if u.IsAdm() {
		statuses.unpublished, err = slices.InSlice(fields, "unpublished")
		if err != nil {
			return "", err
		}
		statuses.deleted, err = slices.InSlice(fields, "deleted")
		if err != nil {
			return "", err
		}
	}

	// We make sure the user picked a status
	if !statuses.deleted && !statuses.unpublished && !statuses.published {
		return "", httperr.NewBadRequest("you need to pick at least 1 valid status")
	}

	if !statuses.deleted {
		// don't want deleted at all
		where = sqlSetAnd(where, "articles.deleted_at IS NULL")
	} else if !statuses.published && !statuses.unpublished && statuses.deleted {
		// _only_ wants deleted articles
		where = sqlSetAnd(where, "articles.deleted_at IS NOT NULL")
	}

	// We set the publish/unPublish flag
	if !statuses.published && statuses.unpublished {
		// only wants unpublished
		where = sqlSetAnd(where, "articles.published_at IS NULL")
	} else if statuses.published && !statuses.unpublished {
		// only wants published
		where = sqlSetAnd(where, "articles.published_at IS NOT NULL")
	}

	return where, nil
}

func sqlSetAnd(clause string, value string) string {
	if clause != "" && value != "" {
		clause += " AND "
	}
	clause += value
	return clause
}

func sqlSetComma(clause string, value string) string {
	if clause != "" && value != "" {
		clause += ", "
	}
	clause += value
	return clause
}
