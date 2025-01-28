package onedrive

import (
	"context"
	http2 "net/http"

	"github.com/bearcatat/onedriver-api/http"
	"github.com/bearcatat/onedriver-api/resources"
)

type core struct {
	client *http.HttpClientWithOauth2
	url    *oneDriveURL
}

func newCore(client *http2.Client) *core {
	return &core{
		client: http.NewHttpClientWithOauth2(client),
		url:    newOneDriveURL(),
	}
}

type Client struct {
	*core
}

func NewClient(client *http2.Client) *Client {
	return &Client{
		core: newCore(client),
	}
}

func (c *Client) GetMyDrive(ctx context.Context) (*Drive, error) {
	req := http.NewJsonRequest(http2.MethodGet, c.url.GetMyDrive(), nil)
	var drive *resources.Drive
	err := c.client.DoWithAuth(ctx, req, &drive)
	if err != nil {
		return nil, err
	}
	return NewDrive(c.core, drive), nil
}
