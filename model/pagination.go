package model

type Pagination struct {
	Total       int `json:"total"`
	PerPage     int `json:"per_page"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
}
