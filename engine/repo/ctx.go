package repo

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	defaultOffset int = 0
	defaultLimit  int = 25
	maxLimit      int = 100
)

// QueryContext determines context values to use
// when generating the DB query.
type QueryCtx struct {
	Limit   int
	Offset  int
	OrderBy string
	Asc     bool
}

// QueryCtxFromGin checks request URL query and extracts
// DB context parameters: limit, offset, sort, desc.
// In case of parsing errors the function will set default values
// to unprocessable fields.
func QueryCtxFromGin(c *gin.Context) QueryCtx {
	ctx := DefaultQueryCtx()
	query := c.Request.URL.Query()
	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, err := strconv.Atoi(queryValue)
			if err != nil || limit > 100 || limit == 0 {
				ctx.Limit = defaultLimit
			} else {
				ctx.Limit = limit
			}
		case "offset":
			// user can attempt to fetch out of bounds
			// records - nothing will be returned back from DB
			offset, err := strconv.Atoi(queryValue)
			if err != nil || offset <= 0 {
				ctx.Offset = defaultOffset
			} else {
				ctx.Offset = offset
			}
		case "sort":
			// no guarantees that the col exists
			if queryValue != "" {
				ctx.OrderBy = queryValue
			}
		case "asc":
			// accept any truthy value
			ctx.Asc = true
		}
	}
	return ctx
}

// DefaultQueryCtx returns default QueryCtx
// with following vals: limit=25, offset=0.
func DefaultQueryCtx() QueryCtx {
	return QueryCtx{
		Limit:  defaultLimit,
		Offset: defaultOffset,
	}
}
