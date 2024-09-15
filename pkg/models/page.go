package models

import (
	"net/http"
	"strconv"
)

// Page defines the pagination struct used for paging results.
type Page struct {
	Limit  uint  `json:"limit"`
	Offset uint  `json:"offset"`
	Total  *uint `json:"total,omitempty"`
}

// DefaultPage returns a default Page struct with a default limit of 15.
func DefaultPage() *Page {
	return &Page{
		Limit:  15,
		Offset: 0,
		Total:  nil,
	}
}

// PageFromRequest constructs a Page from the request. Provides default values is
// values are malformed, missing or invalid.
func PageFromRequest(r *http.Request) *Page {
	page := DefaultPage()
	query := r.URL.Query()

	if query.Has("limit") {
		if limit, err := strconv.Atoi(query.Get("limit")); err == nil {
			if limit <= 0 || limit > 25 {
				page.Limit = 15
			} else {
				page.Limit = uint(limit)
			}
		}
	}

	if query.Has("offset") {
		if offset, err := strconv.Atoi(query.Get("limit")); err == nil {
			if offset > 0 {
				page.Offset = uint(offset)
			}
		}
	}

	return page
}
