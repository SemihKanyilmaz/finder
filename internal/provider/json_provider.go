package provider

import (
	"context"
	"encoding/json"
	"finder/internal/model"
	"fmt"
	client "finder/pkg/http/client"
	"log/slog"
	"net/http"
	"time"
)

type jsonResponse struct {
	Contents []jsonItem `json:"contents"`
}

type jsonItem struct {
	ID          string      `json:"id"`
	Title       string      `json:"title"`
	Type        string      `json:"type"`
	Metrics     jsonMetrics `json:"metrics"`
	PublishedAt string      `json:"published_at"`
}

type jsonMetrics struct {
	Views       int    `json:"views"`
	Likes       int    `json:"likes"`
	Duration    string `json:"duration"`
	ReadingTime int    `json:"reading_time"`
	Reactions   int    `json:"reactions"`
}

type jsonProvider struct {
	name   string
	client *client.Client
}

func NewJSONProvider(name string, c *client.Client) ContentProvider {
	return &jsonProvider{name: name, client: c}
}

func (p *jsonProvider) Fetch(ctx context.Context) ([]model.Content, error) {
	resp, err := p.client.Get(ctx, client.Request{})
	if err != nil {
		return nil, fmt.Errorf("fetch: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}

	var data jsonResponse
	if err := json.Unmarshal(resp.Body, &data); err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	contents := make([]model.Content, 0, len(data.Contents))
	for _, item := range data.Contents {
		c, err := mapJSONItem(item, p.name)
		if err != nil {
			slog.Error("json provider map failed", "id", item.ID, "error", err)
			continue
		}
		contents = append(contents, c)
	}

	return contents, nil
}

func mapJSONItem(item jsonItem, source string) (model.Content, error) {
	publishedAt, err := time.Parse(time.RFC3339, item.PublishedAt)
	if err != nil {
		return model.Content{}, fmt.Errorf("parse time: %w", err)
	}

	contentType := model.ContentType(item.Type)

	return model.Content{
		ID:          item.ID,
		Title:       item.Title,
		Type:        contentType,
		Source:      source,
		PublishedAt: publishedAt,
		Views:       item.Metrics.Views,
		Likes:       item.Metrics.Likes,
		ReadingTime: item.Metrics.ReadingTime,
		Reactions:   item.Metrics.Reactions,
	}, nil
}
