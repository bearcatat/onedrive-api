package http

import (
	"context"
	"fmt"
	"io"
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

func (c *HttpClientWithOauth2) DoWithoutAuth(ctx context.Context, req Request, target interface{}) error {
	return NewHttpClient(&http.Client{}).DoRequestAndParseResponse(ctx, req, target)
}

func (c *HttpClientWithOauth2) DoWithAuth(ctx context.Context, req Request, target interface{}) error {
	return c.oauth2Client.DoRequestAndParseResponse(ctx, req, target)
}

func (c *HttpClientWithOauth2) Download(ctx context.Context, req Request, writer io.Writer) error {
	return c.oauth2Client.Download(ctx, req, writer)
}

type HttpClient struct {
	client *http.Client
}

func NewHttpClient(client *http.Client) *HttpClient {
	return &HttpClient{
		client: client,
	}
}

func (c *HttpClient) DoRequestAndParseResponse(ctx context.Context, req Request, target interface{}) error {
	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	jsonResp := NewJsonResponse(resp)
	return jsonResp.UnMarshal(target)
}

func (c *HttpClient) Download(ctx context.Context, req Request, writer io.Writer) error {
	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	_, err = io.Copy(writer, resp.Body)
	return err
}

func (c *HttpClient) do(ctx context.Context, request Request) (*http.Response, error) {
	httpRequest, err := request.GetHttpRequest()
	if err != nil {
		return nil, err
	}
	httpRequest = httpRequest.WithContext(ctx)
	resp, err := c.client.Do(httpRequest)
	if err != nil {
		return nil, c.handleError(ctx, err)
	}
	return resp, nil
}

func (c *HttpClient) handleError(ctx context.Context, err error) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return err
	}
}
