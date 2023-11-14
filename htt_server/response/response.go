package response

// PageDocument is the form used for API responses from query API calls.
type PageDocument[T any] struct {
	Items       []T `json:"items"`
	Total       int `json:"total"`
	Page        int `json:"page"`
	RowsPerPage int `json:"rows_per_page"`
}

// NewPageDocument constructs a response value for a web paging response.
func NewPageDocument[T any](items []T, total int, page int, rowsPrePage int) PageDocument[T] {
	return PageDocument[T]{
		Items:       items,
		Total:       total,
		Page:        page,
		RowsPerPage: rowsPrePage,
	}
}
