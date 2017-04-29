package httplib

import (
	"net/http"
	"services/config"
	"strconv"
)

var (
	maxPerPage = config.RegisterInt("services.httplib.max_per_page", 100)
	minPerPage = config.RegisterInt("services.httplib.min_per_page", 1)
	perPage    = config.RegisterInt("services.httplib.per_page", 10)
)

// GetPageAndCount return the p and c variable from the request, if not available
// return the default value
func GetPageAndCount(r *http.Request, offset bool) (int, int) {
	p64, err := strconv.ParseInt(r.URL.Query().Get("p"), 10, 0)
	p := int(p64)
	if err != nil || p < 1 {
		p = 1
	}

	c64, err := strconv.ParseInt(r.URL.Query().Get("c"), 10, 0)
	c := int(c64)
	if err != nil || c > *maxPerPage || c < *minPerPage {
		c = *perPage
	}

	if offset {
		// If i need to make it to offset model then do it here
		p = (p - 1) * c
	}

	return p, c
}
