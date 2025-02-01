package onedrive

import (
	"context"
	http2 "net/http"

	"github.com/bearcatat/onedrive-api/http"
	"github.com/bearcatat/onedrive-api/resources"
)

type DriveItem struct {
	*core
	*resources.DriveItem

	drive *resources.Drive
}

func NewDriveItem(c *core, driveItem *resources.DriveItem, drive *resources.Drive) *DriveItem {
	return &DriveItem{
		core:      c,
		DriveItem: driveItem,
		drive:     drive,
	}
}

func (i *DriveItem) CreateFolder(ctx context.Context, folderName string) (*DriveItem, error) {
	var driveItem *resources.DriveItem
	err := i.client.DoWithAuth(ctx, i.createFolderRequest(folderName), &driveItem)
	if err != nil {
		return nil, err
	}
	return NewDriveItem(i.core, driveItem, i.drive), nil
}

func (i *DriveItem) createFolderRequest(folderName string) http.Request {
	url := i.url.CreateFolder(i.drive.Id, i.DriveItem.Id)
	body := resources.NewCreateFolderRequest(folderName)
	return http.NewJsonRequest(http2.MethodPost, url, body)
}

func (i *DriveItem) UploadLargeFile(ctx context.Context, file File) (*DriveItem, error) {
	if file.IsDir() {
		return nil, ErrNotFile
	}
	if file.Size() == 0 {
		return nil, ErrEmptyFile
	}
	uploadSession, err := i.createUploadSession(ctx, file)
	if err != nil {
		return nil, err
	}
	return i.uploadFileToUploadSession(ctx, file, uploadSession)
}

func (i *DriveItem) createUploadSession(ctx context.Context, file File) (*resources.UploadSession, error) {
	var uploadSession *resources.UploadSession
	err := i.client.DoWithAuth(ctx, i.createUploadSessionRequest(file), &uploadSession)
	if err != nil {
		return nil, err
	}
	return uploadSession, nil
}

func (i *DriveItem) createUploadSessionRequest(file File) http.Request {
	url := i.url.UploadSession(i.drive.Id, i.DriveItem.Id, file.Name())
	body := resources.NewUploadSessionRequest()
	return http.NewJsonRequest(http2.MethodPost, url, body)
}

func (i *DriveItem) uploadFileToUploadSession(ctx context.Context, file File, session *resources.UploadSession) (*DriveItem, error) {
	fileForUpload := newFileForUpload(file, session)
	var response *resources.UploadSessionResponse
	for fileForUpload.next() {
		var err error
		response, err = i.uploadFileFragmentToUploadSession(ctx, fileForUpload)
		if err != nil {
			return nil, err
		}
	}
	return NewDriveItem(i.core, &response.DriveItem, i.drive), nil
}

func (i *DriveItem) uploadFileFragmentToUploadSession(ctx context.Context, file *fileForUpload) (*resources.UploadSessionResponse, error) {
	var response *resources.UploadSessionResponse
	req, err := file.getNextRequest()
	if err != nil {
		return nil, err
	}
	err = i.client.Do(ctx, req, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (i *DriveItem) Update(ctx context.Context, update *DriveItem) (*DriveItem, error) {
	var driveItem *resources.DriveItem
	err := i.client.DoWithAuth(ctx, i.updateRequest(update), &driveItem)
	if err != nil {
		return nil, err
	}
	return NewDriveItem(i.core, driveItem, i.drive), nil
}

func (i *DriveItem) updateRequest(item *DriveItem) http.Request {
	url := i.url.Update(i.drive.Id, i.DriveItem.Id)
	return http.NewJsonRequest(http2.MethodPatch, url, item)
}

func (i *DriveItem) Delete(ctx context.Context) error {
	return i.client.DoWithAuth(ctx, i.deleteReqeust(), nil)
}

func (i *DriveItem) deleteReqeust() http.Request {
	url := i.url.Delete(i.drive.Id, i.DriveItem.Id)
	return http.NewJsonRequest(http2.MethodDelete, url, nil)
}

func (i *DriveItem) Copy(ctx context.Context, parentItem *DriveItem) (*AsyncJob, error) {
	var asyncJob *resources.AsyncJob
	err := i.client.DoWithAuth(ctx, i.copyRequest(parentItem), &asyncJob)
	if err != nil {
		return nil, err
	}
	return NewAsyncJob(i.core, asyncJob, i.drive), nil
}

func (i *DriveItem) copyRequest(parentItem *DriveItem) http.Request {
	url := i.url.Copy(i.drive.Id, i.DriveItem.Id)
	return http.NewJsonRequest(http2.MethodPost, url, resources.NewCopyRequest(parentItem.DriveItem, parentItem.drive))
}
