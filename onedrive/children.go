package onedrive

import (
	"context"
	http2 "net/http"
	"net/url"

	"github.com/bearcatat/onedrive-api/http"
	"github.com/bearcatat/onedrive-api/resources"
)

type Children struct {
	core  *core
	raw   *resources.Children
	drive *resources.Drive
	Value []*DriveItem
}

func newChildren(c *core, raw *resources.Children, drive *resources.Drive) *Children {
	items := make([]*DriveItem, 0)
	for _, item := range raw.Value {
		items = append(items, newDriveItem(c, &item, drive))
	}
	return &Children{
		core:  c,
		raw:   raw,
		drive: drive,
		Value: items,
	}
}

func (c *Children) HasNext() bool {
	return c.raw.NextURL != ""
}

func (c *Children) Next(ctx context.Context) (*Children, error) {
	if !c.HasNext() {
		return nil, ErrChildrenNoNext
	}

	var children *resources.Children
	err := c.core.client.DoWithAuth(ctx, c.nextRequest(), &children)
	if err != nil {
		return nil, err
	}
	return newChildren(c.core, children, c.drive), nil

}

func (c *Children) nextRequest() http.Request {
	url, _ := url.Parse(c.raw.NextURL)
	return http.NewJsonRequest(http2.MethodGet, url, nil)
}
