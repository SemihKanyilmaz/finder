package model

import "time"

type ContentType string

const (
	ContentTypeVideo ContentType = "video"
	ContentTypeText  ContentType = "text"
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
	Keyword     string
	ContentType ContentType
	SortBy      string
	Page        int
	PageSize    int
}

type SearchResult struct {
	Items      []Content
	TotalCount int
	Page       int
	PageSize   int
}
