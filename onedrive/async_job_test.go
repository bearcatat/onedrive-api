package onedrive

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/bearcatat/onedriver-api/resources"
)

func setup_async_job() (asyncJob *AsyncJob, mux *http.ServeMux, teardown func()) {
	url, mux, teardown := setup()
	core := newCore(&http.Client{})
	asyncJobURL := url.String() + "async_job"
	asyncJob = NewAsyncJob(core, &resources.AsyncJob{Url: asyncJobURL}, &resources.Drive{Id: "fake_drive_id"})
	asyncJob.url.baseURL = url
	return asyncJob, mux, teardown
}

func TestAsyncJob_GetResource_Finished(t *testing.T) {
	asyncJob, mux, teardown := setup_async_job()
	defer teardown()
	asyncJob.status = &resources.AsyncJobStatus{
		Status:     resources.FINISHED,
		ResourceId: "fake_drive_item_id",
	}

	mux.HandleFunc("/drives/fake_drive_id/items/fake_drive_item_id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		jsonData := readFile(t, "fake_drive_item.json")
		fmt.Fprint(w, string(jsonData))
		w.WriteHeader(http.StatusOK)
	})

	ctx := context.Background()
	item, err := asyncJob.GetResource(ctx)
	if err != nil {
		t.Errorf("AsyncJob.GetResource returned error: %v", err)
	}
	expectedItem := getDataFromFile[*resources.DriveItem](t, "fake_drive_item.json")
	if !reflect.DeepEqual(item.DriveItem, expectedItem) {
		t.Errorf("AsyncJob.GetResource returned %+v, want %+v", item.DriveItem, expectedItem)
	}
}

func TestAsyncJob_GetResource_Finshed_NilStatus(t *testing.T) {
	asyncJob, mux, teardown := setup_async_job()
	defer teardown()

	mux.HandleFunc("/drives/fake_drive_id/items/fake_drive_item_id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		jsonData := readFile(t, "fake_drive_item.json")
		fmt.Fprint(w, string(jsonData))
		w.WriteHeader(http.StatusOK)
	})

	mux.HandleFunc("/async_job", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		jsonData := readFile(t, "fake_finshed_async_job.json")
		fmt.Fprint(w, string(jsonData))
		w.WriteHeader(http.StatusOK)
	})
	ctx := context.Background()
	item, err := asyncJob.GetResource(ctx)
	if err != nil {
		t.Errorf("AsyncJob.GetResource returned error: %v", err)
	}
	expectedItem := getDataFromFile[*resources.DriveItem](t, "fake_drive_item.json")
	if !reflect.DeepEqual(item.DriveItem, expectedItem) {
		t.Errorf("AsyncJob.GetResource returned %+v, want %+v", item.DriveItem, expectedItem)
	}
	expectedAsyncJobStatus := getDataFromFile[*resources.AsyncJobStatus](t, "fake_finshed_async_job.json")
	if !reflect.DeepEqual(asyncJob.status, expectedAsyncJobStatus) {
		t.Errorf("AsyncJob.GetResource returned %+v, want %+v", asyncJob.status, expectedAsyncJobStatus)
	}
}

func TestAsyncJob_GetResource_NotFinished(t *testing.T) {
	asyncJob, _, teardown := setup_async_job()
	defer teardown()
	asyncJob.status = &resources.AsyncJobStatus{
		Status: "inProgress",
	}
	ctx := context.Background()
	_, err := asyncJob.GetResource(ctx)
	if err == nil {
		t.Errorf("AsyncJob.GetResource returned nil, want error")
	}
}

func TestAsyncJob_GetResource_Failed(t *testing.T) {
	asyncJob, mux, teardown := setup_async_job()
	defer teardown()

	mux.HandleFunc("async_job", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		jsonData := readFile(t, "fake_error.json")
		fmt.Fprint(w, string(jsonData))
		w.WriteHeader(http.StatusNotFound)
	})
	ctx := context.Background()
	_, err := asyncJob.GetResource(ctx)
	if err == nil {
		t.Errorf("AsyncJob.GetResource returned nil, want error")
	}
}

func TestAsyncJob_IsFinished_Finished(t *testing.T) {
	asyncJob, _, teardown := setup_async_job()
	defer teardown()
	asyncJob.status = &resources.AsyncJobStatus{
		Status: resources.FINISHED,
	}

	ctx := context.Background()
	isFinished, err := asyncJob.IsFinished(ctx)
	if err != nil {
		t.Errorf("AsyncJob.IsFinished returned error: %v", err)
	}
	if !isFinished {
		t.Errorf("AsyncJob.IsFinished returned false, want true")
	}
}

func TestAsyncJob_IsFinished_InProgress(t *testing.T) {
	asyncJob, mux, teardown := setup_async_job()
	defer teardown()

	mux.HandleFunc("/async_job", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		jsonData := readFile(t, "fake_in_progress_async_job.json")
		fmt.Fprint(w, string(jsonData))
		w.WriteHeader(http.StatusOK)
	})
	ctx := context.Background()
	isFinished, err := asyncJob.IsFinished(ctx)
	if err != nil {
		t.Errorf("AsyncJob.IsFinished returned error: %v", err)
	}
	if isFinished {
		t.Errorf("AsyncJob.IsFinished returned true, want false")
	}
}

func TestAsyncJob_IsFinished_Failed(t *testing.T) {
	asyncJob, mux, teardown := setup_async_job()
	defer teardown()

	mux.HandleFunc("/async_job", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		jsonData := readFile(t, "fake_error.json")
		fmt.Fprint(w, string(jsonData))
	})
	ctx := context.Background()
	_, err := asyncJob.IsFinished(ctx)
	if err == nil {
		t.Errorf("AsyncJob.IsFinished returned nil, want error")
	}
}
