package onedrive

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/bearcatat/onedrive-api/resources"
)

func setup_client() (client *Client, mux *http.ServeMux, teardown func()) {
	url, mux, teardown := setup()
	client = NewClient(&http.Client{})
	client.url.baseURL = url
	return client, mux, teardown
}

func TestClient_GetMyDrive(t *testing.T) {
	client, mux, teardown := setup_client()
	defer teardown()

	mux.HandleFunc("/me/drive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		jsonData := readFile(t, "fake_drive.json")
		fmt.Fprint(w, string(jsonData))
	})

	ctx := context.Background()
	drive, err := client.GetMyDrive(ctx)
	if err != nil {
		t.Errorf("Client.GetMyDrive returned error: %v", err)
	}
	expectedDrive := getDataFromFile[*resources.Drive](t, "fake_drive.json")
	if !reflect.DeepEqual(drive.Drive, expectedDrive) {
		t.Errorf("Client.GetMyDrive returned %+v, want %+v", drive.Drive, expectedDrive)
	}
}

func TestClient_GetMyDrive_Fail(t *testing.T) {
	client, mux, teardown := setup_client()
	defer teardown()

	mux.HandleFunc("/me/drive", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		jsonData := readFile(t, "fake_error.json")
		fmt.Fprint(w, string(jsonData))
		w.WriteHeader(http.StatusNotFound)
	})

	ctx := context.Background()
	_, err := client.GetMyDrive(ctx)
	if err == nil {
		t.Errorf("Client.GetMyDrive should return error")
	}
}
