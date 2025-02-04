package onedrive

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/bearcatat/onedrive-api/resources"
)

func setup_drive_item() (driveItem *DriveItem, mux *http.ServeMux, teardown func()) {
	url, mux, teardown := setup()
	core := newCore(&http.Client{})
	driveItem = newDriveItem(core, &resources.DriveItem{Id: "fake_drive_item_id"}, &resources.Drive{Id: "fake_drive_id"})
	driveItem.url.baseURL = url
	return driveItem, mux, teardown
}

func TestDriveItem_CreateFolder(t *testing.T) {
	driveItem, mux, teardown := setup_drive_item()
	defer teardown()

	mux.HandleFunc("/drives/fake_drive_id/items/fake_drive_item_id/children", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		expectedRequestBody := getDataFromFile[*resources.CreateFolderRequest](t, "fake_create_folder_request_body.json")
		testBody(t, r, expectedRequestBody)

		jsonData := readFile(t, "fake_drive_item.json")
		fmt.Fprint(w, string(jsonData))
	})

	ctx := context.Background()
	item, err := driveItem.CreateFolder(ctx, "test_folder")
	if err != nil {
		t.Errorf("DriveItem.CreateFolder returned error: %v", err)
	}
	expectedItem := getDataFromFile[*resources.DriveItem](t, "fake_drive_item.json")
	if !reflect.DeepEqual(item.DriveItem, expectedItem) {
		t.Errorf("DriveItem.CreateFolder returned %+v, want %+v", item.DriveItem, expectedItem)
	}
}

func TestDriveItem_createUploadSession(t *testing.T) {
	driveItem, mux, teardown := setup_drive_item()
	defer teardown()

	mux.HandleFunc("/drives/fake_drive_id/items/fake_drive_item_id:/fake_file_name:/createUploadSession", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		expectedRequestBody := getDataFromFile[*resources.UploadSessionRequest](t, "fake_create_upload_session_request_body.json")
		testBody(t, r, expectedRequestBody)

		jsonData := readFile(t, "fake_upload_session.json")
		fmt.Fprint(w, string(jsonData))
	})

	ctx := context.Background()
	uploadSession, err := driveItem.createUploadSession(ctx, &fakeFile{})
	if err != nil {
		t.Errorf("DriveItem.createUploadSession returned error: %v", err)
	}
	expectedUploadSession := getDataFromFile[*resources.UploadSession](t, "fake_upload_session.json")
	if !reflect.DeepEqual(uploadSession, expectedUploadSession) {
		t.Errorf("DriveItem.createUploadSession returned %+v, want %+v", uploadSession, expectedUploadSession)
	}
}

func TestDriveItem_UploadLargeFile(t *testing.T) {
	driveItem, mux, teardown := setup_drive_item()
	defer teardown()
	fakeFile := &fakeFile{readTimes: 0}

	mux.HandleFunc("/drives/fake_drive_id/items/fake_drive_item_id:/fake_file_name:/createUploadSession", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		data := getDataFromRequest[*resources.UploadSession](t, r)
		data.UploadURL = driveItem.url.baseURL.String() + "fake_upload_url"
		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Errorf("readTestData failed: %v", err)
		}
		fmt.Fprint(w, string(jsonData))
	})
	mux.HandleFunc("/fake_upload_url", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		testHeader(t, r, "Content-Range", fakeFile.ContentRange())
		testHeader(t, r, "Content-Length", fakeFile.ContentLength())
		if fakeFile.Finished() {
			w.WriteHeader(http.StatusCreated)
			jsonData := readFile(t, "fake_drive_item.json")
			fmt.Fprint(w, string(jsonData))
		} else {
			w.WriteHeader(http.StatusAccepted)
			jsonData := readFile(t, "fake_upload_session_response.json")
			fmt.Fprint(w, string(jsonData))
		}
	})

	item, err := driveItem.UploadLargeFile(context.Background(), fakeFile)
	if err != nil {
		t.Errorf("DriveItem.UploadLargeFile returned error: %v", err)
	}
	expectedItem := getDataFromFile[*resources.DriveItem](t, "fake_drive_item.json")
	if !reflect.DeepEqual(item.DriveItem, expectedItem) {
		t.Errorf("DriveItem.UploadLargeFile returned %+v, want %+v", item.DriveItem, expectedItem)
	}
}

func TestDriveItem_UploadLargeFile_NotFile(t *testing.T) {
	driveItem, _, teardown := setup_drive_item()
	defer teardown()
	fakeDir := &fakeDir{}
	_, err := driveItem.UploadLargeFile(context.Background(), fakeDir)
	if err != ErrNotFile {
		t.Errorf("DriveItem.UploadLargeFile returned %v, want %v", err, ErrNotFile)
	}
}

func TestDriveItem_UploadLargeFile_EmptyFile(t *testing.T) {
	driveItem, _, teardown := setup_drive_item()
	defer teardown()

	emptyFile := &fakeEmptyFile{}
	_, err := driveItem.UploadLargeFile(context.Background(), emptyFile)
	if err != ErrEmptyFile {
		t.Errorf("DriveItem.UploadLargeFile returned %v, want %v", err, ErrEmptyFile)
	}
}

func TestDriveItem_UploadLargeFile_ReadFailed(t *testing.T) {
	driveItem, mux, teardown := setup_drive_item()
	defer teardown()

	mux.HandleFunc("/drives/fake_drive_id/items/fake_drive_item_id:/fake_error_file_name:/createUploadSession", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		data := getDataFromRequest[*resources.UploadSession](t, r)
		data.UploadURL = driveItem.url.baseURL.String() + "fake_upload_url"
		jsonData, err := json.Marshal(data)
		if err != nil {
			t.Errorf("readTestData failed: %v", err)
		}
		fmt.Fprint(w, string(jsonData))
	})

	fakeErrorFile := &fakeErrorFile{}
	_, err := driveItem.UploadLargeFile(context.Background(), fakeErrorFile)
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("DriveItem.UploadLargeFile returned %v, want %v", err, os.ErrNotExist)
	}
}

func TestDiverItem_Update(t *testing.T) {
	driveItem, mux, teardown := setup_drive_item()
	defer teardown()

	mux.HandleFunc("/drives/fake_drive_id/items/fake_drive_item_id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		jsonData := readFile(t, "fake_drive_item.json")
		fmt.Fprint(w, string(jsonData))
	})

	ctx := context.Background()
	item, err := driveItem.Update(ctx, &DriveItem{})
	if err != nil {
		t.Errorf("DriveItem.Update returned error: %v", err)
	}
	expectedItem := getDataFromFile[*resources.DriveItem](t, "fake_drive_item.json")
	if !reflect.DeepEqual(item.DriveItem, expectedItem) {
		t.Errorf("DriveItem.Update returned %+v, want %+v", item.DriveItem, expectedItem)
	}
}

func TestDriveItem_Delete(t *testing.T) {
	driveItem, mux, teardown := setup_drive_item()
	defer teardown()

	mux.HandleFunc("/drives/fake_drive_id/items/fake_drive_item_id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	ctx := context.Background()
	err := driveItem.Delete(ctx)
	if err != nil {
		t.Errorf("DriveItem.Delete returned error: %v", err)
	}
}

func TestDriveItem_Copy(t *testing.T) {
	driveItem, mux, teardown := setup_drive_item()
	defer teardown()

	url := driveItem.url.baseURL.String() + "async_job"
	mux.HandleFunc("/drives/fake_drive_id/items/fake_drive_item_id/copy", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		expectedRequestBody := getDataFromFile[*resources.DriveItem](t, "fake_copy_request_body.json")
		testBody(t, r, expectedRequestBody)

		w.Header().Set("Location", url)
		w.WriteHeader(http.StatusAccepted)
		w.Write(nil)
	})

	ctx := context.Background()
	asyncJob, err := driveItem.Copy(ctx, &DriveItem{
		DriveItem: &resources.DriveItem{
			Id: "fake_parent_drive_item_id",
		},
		drive: driveItem.drive,
	}, "fake_name")
	if err != nil {
		t.Errorf("DriveItem.Copy returned error: %v", err)
	}
	if asyncJob.Url != url {
		t.Errorf("DriveItem.Copy returned %+v, want %+v", asyncJob.Url, url)
	}
}

func TestDriveItem_Move(t *testing.T) {
	driveItem, mux, teardown := setup_drive_item()
	defer teardown()

	mux.HandleFunc("/drives/fake_drive_id/items/fake_drive_item_id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PATCH")
		expectedRequestBody := getDataFromFile[*resources.DriveItem](t, "fake_move_request_body.json")
		testBody(t, r, expectedRequestBody)

		jsonData := readFile(t, "fake_drive_item.json")
		fmt.Fprint(w, string(jsonData))
	})

	ctx := context.Background()
	item, err := driveItem.Move(ctx, &DriveItem{
		DriveItem: &resources.DriveItem{
			Id: "fake_parent_drive_item_id",
		},
		drive: driveItem.drive,
	}, "fake_name")
	if err != nil {
		t.Errorf("DriveItem.Move returned error: %v", err)
	}
	expectedItem := getDataFromFile[*resources.DriveItem](t, "fake_drive_item.json")
	if !reflect.DeepEqual(item.DriveItem, expectedItem) {
		t.Errorf("DriveItem.Move returned %+v, want %+v", item.DriveItem, expectedItem)
	}
}

func TestDriveItem_ListChildren(t *testing.T) {
	driveItem, mux, teardown := setup_drive_item()
	defer teardown()

	mux.HandleFunc("/drives/fake_drive_id/items/fake_drive_item_id/children", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		jsonData := readFile(t, "fake_children.json")
		fmt.Fprint(w, string(jsonData))
	})

	ctx := context.Background()
	children, err := driveItem.ListChildren(ctx)
	if err != nil {
		t.Errorf("DriveItem.ListChildren returned error: %v", err)
	}
	expectedChildren := getDataFromFile[*resources.Children](t, "fake_children.json")
	if !reflect.DeepEqual(children.raw, expectedChildren) {
		t.Errorf("DriveItem.ListChildren returned %+v, want %+v", children, expectedChildren)
	}
}

func TestDriveItem_Download(t *testing.T) {
	driveItem, mux, teardown := setup_drive_item()
	defer teardown()

	downloadURL := driveItem.url.baseURL.String() + "fake_download_url"
	mux.HandleFunc("/drives/fake_drive_id/items/fake_drive_item_id/content", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Location", downloadURL)
		w.WriteHeader(http.StatusFound)
	})

	mux.HandleFunc("/fake_download_url", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		payload := readFile(t, "fake_file.txt")
		n, err := w.Write(payload)
		if err != nil {
			t.Errorf("DriveItem.Download returned error: %v", err)
		}
		t.Log(n)
	})

	ctx := context.Background()
	writer := &bytes.Buffer{}
	err := driveItem.Download(ctx, writer)
	if err != nil {
		t.Errorf("DriveItem.Download returned error: %v", err)
	}

	expected := readFile(t, "fake_file.txt")
	if !bytes.Equal(writer.Bytes(), expected) {
		t.Errorf("DriveItem.Download returned %+v, want %+v", writer.Bytes(), expected)
	}
}

func TestDriveItem_Download_UseDownloadURL(t *testing.T) {
	driveItem, mux, teardown := setup_drive_item()
	defer teardown()

	downloadURL := driveItem.url.baseURL.String() + "fake_download_url"
	driveItem.DownloadURL = downloadURL

	mux.HandleFunc("/fake_download_url", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/octet-stream")
		w.WriteHeader(http.StatusOK)
		payload := readFile(t, "fake_file.txt")
		n, err := w.Write(payload)
		if err != nil {
			t.Errorf("DriveItem.Download returned error: %v", err)
		}
		t.Log(n)
	})

	ctx := context.Background()
	writer := &bytes.Buffer{}
	err := driveItem.Download(ctx, writer)
	if err != nil {
		t.Errorf("DriveItem.Download returned error: %v", err)
	}

	expected := readFile(t, "fake_file.txt")
	if !bytes.Equal(writer.Bytes(), expected) {
		t.Errorf("DriveItem.Download returned %+v, want %+v", writer.Bytes(), expected)
	}
}
