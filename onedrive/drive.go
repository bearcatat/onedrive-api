package onedrive

import (
	"context"
	http2 "net/http"

	"github.com/bearcatat/onedriver-api/http"
	"github.com/bearcatat/onedriver-api/resources"
)

type Drive struct {
	*core

	*resources.Drive
}

func NewDrive(c *core, drive *resources.Drive) *Drive {
	return &Drive{
		core:  c,
		Drive: drive,
	}
}

func (d *Drive) GetByPath(ctx context.Context, path string) (*DriveItem, error) {
	var driverItem *resources.DriveItem
	err := d.client.DoWithAuth(ctx, d.getByPathRequest(path), &driverItem)
	if err != nil {
		return nil, err
	}
	return NewDriveItem(d.core, driverItem, d.Drive), nil
}

func (d *Drive) getByPathRequest(path string) http.Request {
	url := d.url.GetDriveItemByPath(d.Drive.Id, path)
	return http.NewJsonRequest(http2.MethodGet, url, nil)
}

func (d *Drive) Get(ctx context.Context, itemId string) (*DriveItem, error) {
	var driverItem *resources.DriveItem
	err := d.client.DoWithAuth(ctx, d.getRequest(itemId), &driverItem)
	if err != nil {
		return nil, err
	}
	return NewDriveItem(d.core, driverItem, d.Drive), nil
}

func (d *Drive) getRequest(itemId string) http.Request {
	url := d.url.Get(d.Drive.Id, itemId)
	return http.NewJsonRequest(http2.MethodGet, url, nil)
}
