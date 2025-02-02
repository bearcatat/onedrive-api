package onedrive

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	"github.com/bearcatat/onedrive-api/resources"
)

func setup_children() (children *Children, mux *http.ServeMux, teardown func()) {
	url, mux, teardown := setup()
	core := newCore(&http.Client{})
	children = newChildren(core, &resources.Children{}, &resources.Drive{Id: "fake_drive_id"})
	children.raw.NextURL = url.String() + "/next"
	return children, mux, teardown
}

func TestChildren_Next(t *testing.T) {
	children, mux, teardown := setup_children()
	defer teardown()

	mux.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		jsonData := readFile(t, "fake_children.json")
		w.Write(jsonData)
	})

	next, err := children.Next(context.Background())
	if err != nil {
		t.Errorf("Children.Next failed: %v", err)
	}
	expectedChildren := getDataFromFile[*resources.Children](t, "fake_children.json")
	if !reflect.DeepEqual(next.raw, expectedChildren) {
		t.Errorf("Children.Next returned %+v, want %+v", next.raw, expectedChildren)
	}
}

func TestChildren_Next_NoNext(t *testing.T) {
	children, _, teardown := setup_children()
	defer teardown()
	children.raw.NextURL = ""
	_, err := children.Next(context.Background())
	if err != ErrChildrenNoNext {
		t.Errorf("Children.Next returned %+v, want %+v", err, ErrChildrenNoNext)
	}
}

func TestChildren_Next_Fail(t *testing.T) {
	children, mux, teardown := setup_children()
	defer teardown()
	mux.HandleFunc("/next", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		jsonData := readFile(t, "fake_error.json")
		w.Write(jsonData)
	})
	_, err := children.Next(context.Background())
	if err == nil {
		t.Errorf("Children.Next returned no error, want error")
	}
}
