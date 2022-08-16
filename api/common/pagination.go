package common

type PageResponse[T interface{}] struct {
	PageNumber   int32 `json:"page_number,omitempty"`
	ItemsPerPage int32 `json:"items_per_page,omitempty"`
	LastPage     int32 `json:"last_page,omitempty"`
	Data         []T   `json:"data,omitempty"`
}

func NewPageResponse[T interface{}](pageNumber, itemsPerPage, lastPage int32, data []T) *PageResponse[T] {
	return &PageResponse[T]{
		PageNumber:   pageNumber,
		ItemsPerPage: itemsPerPage,
		LastPage:     lastPage,
		Data:         data,
	}
}

func GetPageOffset(pageNumber, itemsPerPage int32) int32 {
	return (pageNumber - 1) * itemsPerPage
}

func GetLastPage(totalElements int64, pageSize int32) int32 {
	return int32(totalElements)/pageSize + 1
}
