package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultPagination(t *testing.T) {
	p := DefaultPagination()
	assert.Equal(t, 1, p.Page)
	assert.Equal(t, 20, p.Limit)
}

func TestPaginationParams_Skip(t *testing.T) {
	tests := []struct {
		name     string
		page     int
		limit    int
		expected int
	}{
		{"first page", 1, 20, 0},
		{"second page", 2, 20, 20},
		{"third page with limit 10", 3, 10, 20},
		{"large page", 5, 50, 200},
		{"page 1 limit 1", 1, 1, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginationParams{Page: tt.page, Limit: tt.limit}
			assert.Equal(t, tt.expected, p.Skip())
		})
	}
}

func TestPaginationParams_Take(t *testing.T) {
	tests := []struct {
		name     string
		limit    int
		expected int
	}{
		{"default limit", 20, 20},
		{"custom limit", 50, 50},
		{"limit 1", 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := PaginationParams{Page: 1, Limit: tt.limit}
			assert.Equal(t, tt.expected, p.Take())
		})
	}
}

func TestBuildPaginationMeta(t *testing.T) {
	tests := []struct {
		name           string
		page           int
		limit          int
		total          int
		expectedPages  int
		expectedTotal  int
	}{
		{"55 items limit 20", 1, 20, 55, 3, 55},
		{"zero items", 1, 20, 0, 0, 0},
		{"exact fit", 1, 20, 20, 1, 20},
		{"one extra", 1, 20, 21, 2, 21},
		{"single item limit 10", 1, 10, 1, 1, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params := PaginationParams{Page: tt.page, Limit: tt.limit}
			meta := BuildPaginationMeta(params, tt.total)

			assert.Equal(t, tt.page, meta.Page)
			assert.Equal(t, tt.limit, meta.Limit)
			assert.Equal(t, tt.expectedTotal, meta.Total)
			assert.Equal(t, tt.expectedPages, meta.TotalPages)
		})
	}
}
