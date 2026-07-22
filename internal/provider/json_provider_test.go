package provider

import (
	"context"
	"finder/internal/model"
	client "finder/pkg/http/client"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONProviderFetch(t *testing.T) {
	body := `{
		"contents": [
			{
				"id": "v1",
				"title": "Go Tutorial",
				"type": "video",
				"metrics": {"views": 15000, "likes": 1200, "duration": "15:30"},
				"published_at": "2024-03-15T10:00:00Z"
			},
			{
				"id": "t1",
				"title": "Go Article",
				"type": "text",
				"metrics": {"reading_time": 8, "reactions": 320},
				"published_at": "2024-03-14T09:00:00Z"
			}
		]
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(body))
	}))
	defer srv.Close()

	c := client.New(client.Config{BaseURL: srv.URL})
	p := NewJSONProvider("json", c)

	items, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}

	if items[0].ID != "v1" || items[0].Title != "Go Tutorial" || items[0].Type != model.ContentTypeVideo {
		t.Errorf("unexpected video item: %+v", items[0])
	}
	if items[0].Views != 15000 || items[0].Likes != 1200 {
		t.Errorf("unexpected video metrics: views=%d likes=%d", items[0].Views, items[0].Likes)
	}
	if items[0].Source != "json" {
		t.Errorf("got source %q, want %q", items[0].Source, "json")
	}

	if items[1].Type != model.ContentTypeText || items[1].ReadingTime != 8 || items[1].Reactions != 320 {
		t.Errorf("unexpected text item: %+v", items[1])
	}
}

func TestJSONProviderFetchNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := client.New(client.Config{BaseURL: srv.URL})
	p := NewJSONProvider("json", c)

	_, err := p.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestJSONProviderFetchInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not json"))
	}))
	defer srv.Close()

	c := client.New(client.Config{BaseURL: srv.URL})
	p := NewJSONProvider("json", c)

	_, err := p.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestJSONProviderFetchInvalidTime(t *testing.T) {
	body := `{
		"contents": [
			{
				"id": "v1",
				"title": "Bad Time",
				"type": "video",
				"metrics": {"views": 100, "likes": 10},
				"published_at": "not-a-date"
			}
		]
	}`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(body))
	}))
	defer srv.Close()

	c := client.New(client.Config{BaseURL: srv.URL})
	p := NewJSONProvider("json", c)

	items, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("got %d items, want 0 (invalid time should be skipped)", len(items))
	}
}

func TestJSONProviderFetchHTTPError(t *testing.T) {
	c := client.New(client.Config{BaseURL: "http://localhost:1"})
	p := NewJSONProvider("json", c)

	_, err := p.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
