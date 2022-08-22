package common

type PageResponse[T interface{}] struct {
	PageNumber   int32 `json:"page_number"`
	ItemsPerPage int32 `json:"items_per_page"`
	LastPage     int32 `json:"last_page"`
	Data         []T   `json:"data"`
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
	div := int32(totalElements) / pageSize
	if totalElements != 0 && int32(totalElements)%pageSize == 0 {
		return div
	}
	return div + 1
}
