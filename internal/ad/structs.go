package ad

import (
	"errors"
	"time"
)

// Pagination type
type Pagination struct {
	Limit          int `json:"limit,omitempty"`
	TotalCount     int `json:"total_count,omitempty"`
	CurrentPage    int `json:"current_page,omitempty"`
	TotalPageCount int `json:"total_page_count,omitempty"`
}

type Full struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	Price       int       `json:"price"`
	PhotoLinks  []string  `json:"photo_links,omitempty"`
	PhotoMain   string    `json:"photo_main"`
}

type Short struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Price       int       `json:"price"`
	PhotoMain   string    `json:"photo_main"`
}

// AdResponse type
type AdResponse struct {
	Success bool  `json:"success"`
	Ad      *Full `json:"ad,omitempty"`
}

// AdsResponse type
type AdsResponse struct {
	Success    bool        `json:"success"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Ads        []Short      `json:"ads,omitempty"`
}

// AdResponse type
type AdCreateResponse struct {
	Success bool `json:"success"`
	Id      int  `json:"id"`
}

// CreateAdRequest represents an ad creation request.
type CreateAdRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Price       int      `json:"price"`
	PhotoLinks  []string `json:"photo_links"`
}

var (
	ErrIsEmpty = errors.New("is empty")
	ErrMaxLen = errors.New("field length exceeded")
)

// Validate validates the CreateAlbumRequest fields.
func (m CreateAdRequest) Validate() map[string]error {
	fieldErrs := make(map[string]error)
	if m.Title == "" {
		fieldErrs["title"] = ErrIsEmpty
	}
	if len(m.Title) > 200 {
		fieldErrs["title"] = ErrMaxLen
	}
	if len(m.Description) >1000 {
		fieldErrs["description"] = ErrMaxLen
	}
	if len(m.PhotoLinks) >3  {
		fieldErrs["photo_links"] = ErrMaxLen
	}
	if m.Price < 1   {
		fieldErrs["price"] = ErrIsEmpty
	}

	if len(fieldErrs) == 0 {
		return nil
	}
	return fieldErrs
}
