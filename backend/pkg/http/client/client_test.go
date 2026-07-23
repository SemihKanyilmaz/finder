package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Get_OK(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL})
	resp, err := c.Get(context.Background(), Request{Path: "/"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusOK)
	}
	if string(resp.Body) != `{"ok":true}` {
		t.Errorf("unexpected body: %s", resp.Body)
	}
}

func TestClient_Get_WithQueryAndHeader(t *testing.T) {
	var gotQuery, gotHeader string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotQuery = r.URL.Query().Get("key")
		gotHeader = r.Header.Get("X-Custom")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL})
	_, err := c.Get(context.Background(), Request{
		Path:   "/",
		Query:  map[string]string{"key": "value"},
		Header: map[string]string{"X-Custom": "header-val"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotQuery != "value" {
		t.Errorf("query not forwarded, got %q", gotQuery)
	}
	if gotHeader != "header-val" {
		t.Errorf("header not forwarded, got %q", gotHeader)
	}
}

func TestClient_Get_ConnectionError(t *testing.T) {
	c := New(Config{BaseURL: "http://localhost:1"})
	_, err := c.Get(context.Background(), Request{Path: "/"})
	if err == nil {
		t.Fatal("expected error for unreachable server")
	}
}

func TestClient_Get_InvalidURL(t *testing.T) {
	c := New(Config{BaseURL: "://invalid"})
	_, err := c.Get(context.Background(), Request{})
	if err == nil {
		t.Fatal("expected error for invalid URL")
	}
}

func TestClient_Get_ReadBodyError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hj, ok := w.(http.Hijacker)
		if !ok {
			t.Error("server does not support hijacking")
			return
		}
		conn, _, _ := hj.Hijack()
		// send headers claiming 100 bytes but close immediately
		conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Length: 100\r\n\r\n"))
		conn.Close()
	}))
	defer srv.Close()

	c := New(Config{BaseURL: srv.URL})
	_, err := c.Get(context.Background(), Request{Path: "/"})
	if err == nil {
		t.Fatal("expected error for truncated body")
	}
}


func TestClient_DefaultTimeout(t *testing.T) {
	c := New(Config{BaseURL: "http://example.com"})
	if c.http.Timeout != 10*time.Second {
		t.Errorf("expected 10s default timeout, got %v", c.http.Timeout)
	}
}

func TestClient_CustomTimeout(t *testing.T) {
	c := New(Config{BaseURL: "http://example.com", Timeout: 5 * time.Second})
	if c.http.Timeout != 5*time.Second {
		t.Errorf("expected 5s timeout, got %v", c.http.Timeout)
	}
}
