package entities

import "math"

// PaginationParams holds validated pagination parameters
type PaginationParams struct {
	Page  int
	Limit int
}

// DefaultPagination returns default pagination params
func DefaultPagination() PaginationParams {
	return PaginationParams{Page: 1, Limit: 20}
}

// Skip returns the offset for database queries
func (p PaginationParams) Skip() int {
	return (p.Page - 1) * p.Limit
}

// Take returns the limit for database queries
func (p PaginationParams) Take() int {
	return p.Limit
}

// PaginationMeta holds pagination metadata for responses
type PaginationMeta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"totalPages"`
}

// PaginatedResult holds a paginated list of items
type PaginatedResult[T any] struct {
	Data []T            `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// BuildPaginationMeta creates pagination metadata from params and total count
func BuildPaginationMeta(params PaginationParams, total int) PaginationMeta {
	totalPages := 0
	if params.Limit > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(params.Limit)))
	}
	return PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
	}
}
