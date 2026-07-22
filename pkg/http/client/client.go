package client

import (
	"context"
	"fmt"

	"github.com/valyala/fasthttp"
)

type Config struct {
	BaseURL string
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
	c       *fasthttp.Client
	baseURL string
}

func New(cfg Config) *Client {
	return &Client{
		c:       &fasthttp.Client{},
		baseURL: cfg.BaseURL,
	}
}

func (cl *Client) Get(ctx context.Context, r Request) (*Response, error) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodGet)
	req.SetRequestURI(cl.baseURL + r.Path)

	for k, v := range r.Query {
		req.URI().QueryArgs().Set(k, v)
	}
	for k, v := range r.Header {
		req.Header.Set(k, v)
	}

	var err error
	if d, ok := ctx.Deadline(); ok {
		err = cl.c.DoDeadline(req, resp, d)
	} else {
		err = cl.c.Do(req, resp)
	}
	if err != nil {
		return nil, fmt.Errorf("get %s: %w", cl.baseURL+r.Path, err)
	}

	return &Response{
		StatusCode: resp.StatusCode(),
		Body:       append([]byte(nil), resp.Body()...),
	}, nil
}
