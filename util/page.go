package util

func NormPage(page int) int {
	if page <= 0 {
		page = 1
	}
	return page
}

func NormPageSize(pageSize int) int {
	if pageSize <= 0 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return pageSize
}
