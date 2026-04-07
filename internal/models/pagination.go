package models

import "fmt"

type Pagination struct {
	Page  int `json:"page"`
	Items int `json:"items"`
	Count int `json:"count"`
	Pages int `json:"pages"`
}

func (p Pagination) Summary() string {
	return fmt.Sprintf("Page %d of %d (%d total)", p.Page, p.Pages, p.Count)
}

func (p Pagination) HasNext() bool {
	return p.Page < p.Pages
}

func (p Pagination) HasPrev() bool {
	return p.Page > 1
}
