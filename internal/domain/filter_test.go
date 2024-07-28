package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortOrder_Mongo(t *testing.T) {
	tests := []struct {
		name string
		s    SortOrder
		want int
	}{
		{
			name: "asc",
			s:    SortOrderAsc,
			want: 1,
		},
		{
			name: "desc",
			s:    SortOrderDesc,
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.Mongo(); got != tt.want {
				t.Errorf("SortOrder.Mongo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_Limit(t *testing.T) {
	tests := []struct {
		name string
		f    Filter
		want int32
	}{
		{
			name: "limit",
			f:    Filter{PageSize: 10},
			want: 10,
		},
		{
			name: "limit 0",
			f:    Filter{PageSize: 0},
			want: 0,
		},
		{
			name: "limit 100",
			f:    Filter{PageSize: 100},
			want: 100,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Limit(); got != tt.want {
				t.Errorf("Filter.Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFilter_Offset(t *testing.T) {
	tests := []struct {
		name string
		f    Filter
		want int32
	}{
		{
			name: "offset",
			f:    Filter{Page: 2, PageSize: 10},
			want: 10,
		},
		{
			name: "offset 0",
			f:    Filter{Page: 1, PageSize: 0},
			want: 0,
		},
		{
			name: "offset 100",
			f:    Filter{Page: 11, PageSize: 10},
			want: 100,
		},
		{
			name: "offset 80",
			f:    Filter{Page: 5, PageSize: 20},
			want: 80,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.Offset(); got != tt.want {
				t.Errorf("Filter.Offset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWithPageSize(t *testing.T) {
	t.Run("valid page size", func(t *testing.T) {
		filter := &Filter{FilterMap: make(map[string]bool)}
		err := WithPageSize(10)(filter)
		assert.NoError(t, err)
		assert.Equal(t, int32(10), filter.PageSize)
		assert.True(t, filter.FilterMap["page_size"])
	})

	t.Run("invalid page size", func(t *testing.T) {
		filter := &Filter{FilterMap: make(map[string]bool)}
		err := WithPageSize(0)(filter)
		assert.NoError(t, err)
		assert.Equal(t, int32(25), filter.PageSize)
		assert.True(t, filter.FilterMap["page_size"])
	})
}

func TestWithSortBy(t *testing.T) {
	t.Run("valid sort by", func(t *testing.T) {
		filter := &Filter{FilterMap: make(map[string]bool)}
		sortBy := SortBy("name")
		err := WithSortBy(sortBy)(filter)
		assert.NoError(t, err)
		assert.Equal(t, sortBy, filter.SortBy)
		assert.True(t, filter.FilterMap["sort_by"])
	})
}

func TestWithSortOrder(t *testing.T) {
	t.Run("valid sort order", func(t *testing.T) {
		filter := &Filter{FilterMap: make(map[string]bool)}
		sortOrder := SortOrder("asc")
		err := WithSortOrder(sortOrder)(filter)
		assert.NoError(t, err)
		assert.Equal(t, sortOrder, filter.SortOrder)
		assert.True(t, filter.FilterMap["sort_order"])
	})

	t.Run("default sort order", func(t *testing.T) {
		filter := &Filter{FilterMap: make(map[string]bool)}
		err := WithSortOrder("")(filter)
		assert.NoError(t, err)
		assert.Equal(t, SortOrderDesc, filter.SortOrder)
		assert.True(t, filter.FilterMap["sort_order"])
	})
}
