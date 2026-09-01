package dto

type PaginatedDTO struct {
	Limit  int    `json:"limit"`
	Offset int    `json:"offset"`
	Sort   string `json:"sort"`
}

type PaginatedResponse[T any] struct {
	Total int64 `json:"total"` // Total number of items available
	Items []*T  `json:"items"` // List of items in the current page
}

func NewPaginatedResponse[T any](total int64, items []*T) *PaginatedResponse[T] {
	return &PaginatedResponse[T]{
		Total: total,
		Items: items,
	}
}
