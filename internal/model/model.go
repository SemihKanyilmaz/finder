package model

import "time"

type ContentType string

const (
	ContentTypeVideo   ContentType = "video"
	ContentTypeArticle ContentType = "article"
)

type Content struct {
	ID          string
	Title       string
	Type        ContentType
	Source      string
	URL         string
	PublishedAt time.Time
	Views       int
	Likes       int
	ReadingTime int
	Reactions   int
	Score       float64
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
	Page       int       `json:"page"`
	PageSize   int       `json:"pageSize"`
}
