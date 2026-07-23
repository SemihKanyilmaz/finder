package provider

import (
	"context"
	"finder/internal/model"
	client "finder/pkg/http/client"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestXMLProviderFetch(t *testing.T) {
	body := `<?xml version="1.0" encoding="UTF-8"?>
	<feed>
		<items>
			<item>
				<id>v1</id>
				<headline>Docker Intro</headline>
				<type>video</type>
				<stats>
					<views>22000</views>
					<likes>1800</likes>
					<duration>25:15</duration>
				</stats>
				<publication_date>2024-03-15</publication_date>
			</item>
			<item>
				<id>a1</id>
				<headline>Clean Architecture</headline>
				<type>article</type>
				<stats>
					<reading_time>8</reading_time>
					<reactions>450</reactions>
					<comments>25</comments>
				</stats>
				<publication_date>2024-03-14</publication_date>
			</item>
		</items>
	</feed>`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(body))
	}))
	defer srv.Close()

	c := client.New(client.Config{BaseURL: srv.URL})
	p := NewXMLProvider("xml", c)

	items, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Fatalf("got %d items, want 2", len(items))
	}

	if items[0].ID != "v1" || items[0].Title != "Docker Intro" || items[0].Type != model.ContentTypeVideo {
		t.Errorf("unexpected video item: %+v", items[0])
	}
	if items[0].Views != 22000 || items[0].Likes != 1800 {
		t.Errorf("unexpected video stats: views=%d likes=%d", items[0].Views, items[0].Likes)
	}
	if items[0].Source != "xml" {
		t.Errorf("got source %q, want %q", items[0].Source, "xml")
	}

	if items[1].ID != "a1" || items[1].Type != model.ContentTypeArticle {
		t.Errorf("article should map to text type: %+v", items[1])
	}
	if items[1].Title != "Clean Architecture" {
		t.Errorf("headline should map to title: got %q", items[1].Title)
	}
	if items[1].ReadingTime != 8 || items[1].Reactions != 450 {
		t.Errorf("unexpected article stats: reading_time=%d reactions=%d", items[1].ReadingTime, items[1].Reactions)
	}
}

func TestXMLProviderFetchNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	c := client.New(client.Config{BaseURL: srv.URL})
	p := NewXMLProvider("xml", c)

	_, err := p.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestXMLProviderFetchInvalidXML(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("not xml"))
	}))
	defer srv.Close()

	c := client.New(client.Config{BaseURL: srv.URL})
	p := NewXMLProvider("xml", c)

	_, err := p.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error for invalid XML")
	}
}

func TestXMLProviderFetchInvalidDate(t *testing.T) {
	body := `<?xml version="1.0" encoding="UTF-8"?>
	<feed>
		<items>
			<item>
				<id>v1</id>
				<headline>Bad Date</headline>
				<type>video</type>
				<stats><views>100</views><likes>10</likes></stats>
				<publication_date>not-a-date</publication_date>
			</item>
		</items>
	</feed>`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(body))
	}))
	defer srv.Close()

	c := client.New(client.Config{BaseURL: srv.URL})
	p := NewXMLProvider("xml", c)

	items, err := p.Fetch(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("got %d items, want 0 (invalid date should be skipped)", len(items))
	}
}

func TestXMLProviderFetchHTTPError(t *testing.T) {
	c := client.New(client.Config{BaseURL: "http://localhost:1"})
	p := NewXMLProvider("xml", c)

	_, err := p.Fetch(context.Background())
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}
