package model

import (
	"time"
)

type ContentType string

const (
	ContentTypeVideo   ContentType = "video"
	ContentTypeArticle ContentType = "article"
)

type Content struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Type        ContentType `json:"type"`
	Source      string      `json:"source"`
	PublishedAt time.Time   `json:"publishedAt"`
	Score       float64     `json:"score"`
	Views       int         `json:"views"`
	Likes       int         `json:"likes"`
	ReadingTime int         `json:"readingTime"`
	Reactions   int         `json:"reactions"`
}

type SearchParams struct {
	Keyword     string      `query:"q"`
	ContentType ContentType `query:"type"`
	SortBy      string      `query:"sortBy"`
	Page        int         `query:"page"`
	PageSize    int         `query:"pageSize"`
}

type SearchResult struct {
	Items      []Content `json:"items"`
	TotalCount int       `json:"totalCount"`
	TotalPages int       `json:"totalPages"`
	Page       int       `json:"page"`
	PageSize   int       `json:"pageSize"`
}
