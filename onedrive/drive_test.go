package onedrive

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/bearcatat/onedrive-api/resources"
)

func setup_drive() (drive *Drive, mux *http.ServeMux, teardown func()) {
	url, mux, teardown := setup()
	core := newCore(&http.Client{})
	drive = NewDrive(core, &resources.Drive{Id: "fake_drive_id"})
	drive.url.baseURL = url
	return drive, mux, teardown
}

func TestDrive_GetByPath_Success(t *testing.T) {
	drive, mux, teardown := setup_drive()
	defer teardown()

	mux.HandleFunc("/drives/fake_drive_id/root:/path/to/item", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		jsonData := readFile(t, "fake_drive_item.json")
		fmt.Fprint(w, string(jsonData))
	})

	ctx := context.Background()
	item, err := drive.GetByPath(ctx, "path/to/item")
	if err != nil {
		t.Errorf("Drive.GetByPath returned error: %v", err)
	}
	expectedItem := getDataFromFile[*resources.DriveItem](t, "fake_drive_item.json")
	if !reflect.DeepEqual(item.DriveItem, expectedItem) {
		t.Errorf("Drive.GetByPath returned %+v, want %+v", item.DriveItem, expectedItem)
	}
}

func TestDrive_GetByPath_Fail(t *testing.T) {
	drive, mux, teardown := setup_drive()
	defer teardown()

	mux.HandleFunc("/drives/fake_drive_id/root:/path/to/item", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
		jsonData := readFile(t, "fake_error.json")
		fmt.Fprint(w, string(jsonData))
	})

	ctx := context.Background()
	item, err := drive.GetByPath(ctx, "path/to/item")
	if err == nil {
		t.Errorf("Drive.GetByPath should return error")
	}
	if item != nil {
		t.Errorf("Drive.GetByPath should return nil")
	}
}

func TestDrive_Get(t *testing.T) {
	drive, mux, teardown := setup_drive()
	defer teardown()

	mux.HandleFunc("/drives/fake_drive_id/items/fake_item_id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		jsonData := readFile(t, "fake_drive_item.json")
		fmt.Fprint(w, string(jsonData))
	})

	ctx := context.Background()
	item, err := drive.Get(ctx, "fake_item_id")
	if err != nil {
		t.Errorf("Drive.Get returned error: %v", err)
	}
	expectedItem := getDataFromFile[*resources.DriveItem](t, "fake_drive_item.json")
	if !reflect.DeepEqual(item.DriveItem, expectedItem) {
		t.Errorf("Drive.Get returned %+v, want %+v", item.DriveItem, expectedItem)
	}
}

func TestDrive_Get_Fail(t *testing.T) {
	drive, mux, teardown := setup_drive()
	defer teardown()

	mux.HandleFunc("/drives/fake_drive_id/items/fake_item_id", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.WriteHeader(http.StatusNotFound)
		jsonData := readFile(t, "fake_error.json")
		fmt.Fprint(w, string(jsonData))
	})

	ctx := context.Background()
	item, err := drive.Get(ctx, "fake_item_id")
	if err == nil {
		t.Errorf("Drive.Get should return error")
	}
	if item != nil {
		t.Errorf("Drive.Get should return nil")
	}
}
