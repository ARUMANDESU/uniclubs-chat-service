package domain

type SortOrder string
type SortBy string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

func (s SortOrder) Mongo() int {
	if s == SortOrderAsc {
		return 1
	}
	return -1
}

const (
	SortByCreatedAt SortBy = "created_at"
	SortByUpdatedAt SortBy = "updated_at"
)

// made just for fun

// FilterConfiguration is a function that configures a Filter
type FilterConfiguration func(filter *Filter) error

type Filter struct {
	Page      int32
	PageSize  int32
	SortBy    SortBy
	SortOrder SortOrder
	FilterMap map[string]bool
}

func (f Filter) Limit() int32 {
	return f.PageSize
}
func (f Filter) Offset() int32 {
	return (f.Page - 1) * f.PageSize
}

func NewFilter(cfgs ...FilterConfiguration) (*Filter, error) {
	filter := &Filter{
		Page:      1,
		PageSize:  10,
		SortBy:    SortByCreatedAt,
		SortOrder: SortOrderDesc,
		FilterMap: make(map[string]bool),
	}

	for _, cfg := range cfgs {
		err := cfg(filter)
		if err != nil {
			return nil, err
		}
	}

	return filter, nil
}

func WithPage(page int32) FilterConfiguration {
	if page < 1 {
		page = 1
	}

	return func(filter *Filter) error {
		filter.Page = page
		filter.FilterMap["page"] = true
		return nil
	}
}

func WithPageSize(pageSize int32) FilterConfiguration {
	if pageSize < 1 {
		pageSize = 1
	}

	return func(filter *Filter) error {
		filter.PageSize = pageSize
		filter.FilterMap["page_size"] = true
		return nil
	}
}

func WithSortBy(sortBy SortBy) FilterConfiguration {
	if sortBy == "" {
		sortBy = SortByCreatedAt
	}
	return func(filter *Filter) error {
		filter.SortBy = sortBy
		filter.FilterMap["sort_by"] = true
		return nil
	}
}

func WithSortOrder(sortOrder SortOrder) FilterConfiguration {
	if sortOrder == "" {
		sortOrder = SortOrderDesc
	}
	return func(filter *Filter) error {
		filter.SortOrder = sortOrder
		filter.FilterMap["sort_order"] = true
		return nil
	}
}
