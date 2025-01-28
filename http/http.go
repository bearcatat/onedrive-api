package http

import (
	"context"
	"net/http"
)

type HttpClientWithOauth2 struct {
	oauth2Client *HttpClient
}

func NewHttpClientWithOauth2(client *http.Client) *HttpClientWithOauth2 {
	return &HttpClientWithOauth2{
		oauth2Client: NewHttpClient(client),
	}
}

func (c *HttpClientWithOauth2) Do(ctx context.Context, req Request, target interface{}) error {
	return NewHttpClient(&http.Client{}).Do(ctx, req, target)
}

func (c *HttpClientWithOauth2) DoWithAuth(ctx context.Context, req Request, target interface{}) error {
	return c.oauth2Client.Do(ctx, req, target)
}

type HttpClient struct {
	client *http.Client
}

func NewHttpClient(client *http.Client) *HttpClient {
	return &HttpClient{
		client: client,
	}
}

func (c *HttpClient) Do(ctx context.Context, req Request, target interface{}) error {
	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	return resp.UnMarshal(target)
}

func (c *HttpClient) do(ctx context.Context, request Request) (*HttpResponse, error) {
	httpRequest, err := request.GetHttpRequest()
	if err != nil {
		return nil, err
	}
	httpRequest = httpRequest.WithContext(ctx)
	resp, err := c.client.Do(httpRequest)
	if err != nil {
		return nil, c.handleError(ctx, err)
	}
	return NewResponse(resp), nil
}

func (c *HttpClient) handleError(ctx context.Context, err error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return err
	}
}
