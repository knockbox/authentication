package models

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
