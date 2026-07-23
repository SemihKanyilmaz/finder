package client

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Config struct {
	BaseURL string
	Timeout time.Duration
}

type Request struct {
	Path   string
	Query  map[string]string
	Header map[string]string
}

type Response struct {
	StatusCode int
	Body       []byte
}

type Client struct {
	http    *http.Client
	baseURL string
}

func New(cfg Config) *Client {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 10 * time.Second
	}
	return &Client{
		http:    &http.Client{Timeout: timeout},
		baseURL: cfg.BaseURL,
	}
}

func (cl *Client) Get(ctx context.Context, r Request) (*Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cl.baseURL+r.Path, nil)
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}

	q := req.URL.Query()
	for k, v := range r.Query {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	for k, v := range r.Header {
		req.Header.Set(k, v)
	}

	resp, err := cl.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", req.URL.String(), err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       body,
	}, nil
}
