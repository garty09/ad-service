package pagination

import "strconv"

var DefaultPageSize = 10

type Pages struct {
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	PageCount  int         `json:"page_count"`
	TotalCount int         `json:"total_count"`
}

func New(page, total int) *Pages {
	peerPage := DefaultPageSize
	pageCount := -1
	if total >= 0 {
		pageCount = (total + peerPage - 1) / peerPage
		if page > pageCount {
			page = pageCount
		}
	}
	if page < 1 {
		page = 1
	}

	return &Pages{
		Page:       page,
		PerPage:    peerPage,
		TotalCount: total,
		PageCount:  pageCount,
	}
}

func (p *Pages) Offset() int {
	return (p.Page - 1) * p.PerPage
}

func ParseInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}