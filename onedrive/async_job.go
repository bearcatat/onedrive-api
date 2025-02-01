package onedrive

import (
	"context"
	http2 "net/http"
	"net/url"

	"github.com/bearcatat/onedrive-api/http"
	"github.com/bearcatat/onedrive-api/resources"
)

type AsyncJob struct {
	*core
	*resources.AsyncJob

	status *resources.AsyncJobStatus
	drive  *resources.Drive
}

func NewAsyncJob(c *core, asyncJob *resources.AsyncJob, drive *resources.Drive) *AsyncJob {
	return &AsyncJob{
		core:     c,
		AsyncJob: asyncJob,
		drive:    drive,
	}
}

func (a *AsyncJob) GetResource(ctx context.Context) (*DriveItem, error) {
	if a.status == nil {
		err := a.getStatus(ctx)
		if err != nil {
			return nil, err
		}
	}
	if a.status.Status != resources.FINISHED {
		return nil, ErrNotFinished
	}
	drive := NewDrive(a.core, a.drive)
	return drive.Get(ctx, a.status.ResourceId)
}

func (a *AsyncJob) IsFinished(ctx context.Context) (bool, error) {
	if a.status != nil && a.status.Status == resources.FINISHED {
		return true, nil
	}
	err := a.getStatus(ctx)
	if err != nil {
		return false, err
	}
	return a.status.Status == resources.FINISHED, nil
}

func (a *AsyncJob) getStatus(ctx context.Context) error {
	err := a.client.DoWithAuth(ctx, a.getStatusRequest(), &a.status)
	if err != nil {
		return err
	}
	return nil
}

func (a *AsyncJob) getStatusRequest() http.Request {
	url, _ := url.Parse(a.AsyncJob.Url)
	return http.NewJsonRequest(http2.MethodGet, url, nil)
}
