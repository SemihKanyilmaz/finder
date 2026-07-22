package provider

import (
	"context"
	"encoding/xml"
	"finder/internal/model"
	"fmt"
	client "finder/pkg/http/client"
	"log/slog"
	"net/http"
	"time"
)

type xmlFeed struct {
	XMLName xml.Name  `xml:"feed"`
	Items   []xmlItem `xml:"items>item"`
}

type xmlItem struct {
	ID              string   `xml:"id"`
	Headline        string   `xml:"headline"`
	Type            string   `xml:"type"`
	Stats           xmlStats `xml:"stats"`
	PublicationDate string   `xml:"publication_date"`
}

type xmlStats struct {
	Views       int    `xml:"views"`
	Likes       int    `xml:"likes"`
	Duration    string `xml:"duration"`
	ReadingTime int    `xml:"reading_time"`
	Reactions   int    `xml:"reactions"`
}

type xmlProvider struct {
	name   string
	client *client.Client
}

func NewXMLProvider(name string, c *client.Client) ContentProvider {
	return &xmlProvider{name: name, client: c}
}

func (p *xmlProvider) Fetch(ctx context.Context) ([]model.Content, error) {
	resp, err := p.client.Get(ctx, client.Request{})
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var feed xmlFeed
	if err := xml.Unmarshal(resp.Body, &feed); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	contents := make([]model.Content, 0, len(feed.Items))
	for _, item := range feed.Items {
		c, err := mapXMLItem(item, p.name)
		if err != nil {
			slog.Error("xml provider map failed", "id", item.ID, "error", err)
			continue
		}
		contents = append(contents, c)
	}

	return contents, nil
}

func mapXMLItem(item xmlItem, source string) (model.Content, error) {
	publishedAt, err := time.Parse("2006-01-02", item.PublicationDate)
	if err != nil {
		return model.Content{}, fmt.Errorf("parse time: %w", err)
	}

	contentType := model.ContentType(item.Type)
	if contentType == "article" {
		contentType = model.ContentTypeText
	}

	return model.Content{
		ID:          item.ID,
		Title:       item.Headline,
		Type:        contentType,
		Source:      source,
		PublishedAt: publishedAt,
		Views:       item.Stats.Views,
		Likes:       item.Stats.Likes,
		ReadingTime: item.Stats.ReadingTime,
		Reactions:   item.Stats.Reactions,
	}, nil
}
