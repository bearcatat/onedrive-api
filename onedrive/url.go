package onedrive

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	baseURL = "https://graph.microsoft.com/v1.0"
)

// TODO:
//  - Download
//  - List shared files
//  - Recent files
//  - Search
//  - Upload
//  - Query String Parameters

type oneDriveURL struct {
	baseURL *url.URL
}

func newOneDriveURL() *oneDriveURL {
	baseURL, _ := url.Parse(baseURL)
	return &oneDriveURL{
		baseURL: baseURL,
	}
}

func (u *oneDriveURL) GetMyDrive() *url.URL {
	url := u.baseURL.JoinPath("/me/drive")
	return url
}

func (u *oneDriveURL) GetDriveItemByPath(driveId, path string) *url.URL {
	relativePath := fmt.Sprintf("/drives/%s/root", driveId)
	if path != "" {
		path = url.PathEscape(path)
		relativePath = fmt.Sprintf("/drives/%s/root:/%s", driveId, path)
	}
	url := u.baseURL.JoinPath(relativePath)
	return url
}

// GET /drives/{drive-id}/items/{item-id}
func (u *oneDriveURL) Get(driveId, itemId string) *url.URL {
	relativePath := fmt.Sprintf("/drives/%s/items/%s", driveId, itemId)
	url := u.baseURL.JoinPath(relativePath)
	return url
}

func (u *oneDriveURL) CreateFolder(driveId, itemId string) *url.URL {
	relativePath := fmt.Sprintf("/drives/%s/items/%s/children", driveId, itemId)
	url := u.baseURL.JoinPath(relativePath)
	return url
}

func (u *oneDriveURL) UploadSession(driveId, itemId, fileName string) *url.URL {
	fileName = strings.ToValidUTF8(fileName, "x")
	relativePath := fmt.Sprintf("/drives/%s/items/%s:/%s:/createUploadSession", driveId, itemId, fileName)
	url := u.baseURL.JoinPath(relativePath)
	return url
}

func (u *oneDriveURL) Update(driverId, itemId string) *url.URL {
	relativePath := fmt.Sprintf("/drives/%s/items/%s", driverId, itemId)
	return u.baseURL.JoinPath(relativePath)
}

// DELETE /drives/{drive-id}/items/{item-id}
func (u *oneDriveURL) Delete(driverId, itemId string) *url.URL {
	relativePath := fmt.Sprintf("/drives/%s/items/%s", driverId, itemId)
	return u.baseURL.JoinPath(relativePath)
}

// POST /drives/{driveId}/items/{itemId}/copy
func (u *oneDriveURL) Copy(driverId, itemId string) *url.URL {
	relativePath := fmt.Sprintf("/drives/%s/items/%s/copy", driverId, itemId)
	return u.baseURL.JoinPath(relativePath)
}

// PATCH /drives/{drive-id}/items/{item-id}
func (u *oneDriveURL) Move(driverId, itemId string) *url.URL {
	relativePath := fmt.Sprintf("/drives/%s/items/%s", driverId, itemId)
	return u.baseURL.JoinPath(relativePath)
}

// GET /drives/{drive-id}/items/{item-id}/children
func (u *oneDriveURL) ListChildren(driverId, itemId string) *url.URL {
	relativePath := fmt.Sprintf("/drives/%s/items/%s/children", driverId, itemId)
	return u.baseURL.JoinPath(relativePath)
}
