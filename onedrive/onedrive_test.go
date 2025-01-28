package onedrive

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"testing"
)

const (
	// baseURLPath is a non-empty Client.BaseURL path to use during tests,
	// to ensure relative URLs are used for all endpoints.
	baseURLPath         = "/api"
	baseOneDriveURLPath = "/test-onedrive-api"
)

func setup() (url *url.URL, mux *http.ServeMux, teardown func()) {
	// mux is the HTTP request multiplexer used with the test server.
	mux = http.NewServeMux()

	// Ensure that tests catch mistakes where the endpoint URL is specified as absolute rather than relative.
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.Handle(baseOneDriveURLPath+"/", http.StripPrefix(baseOneDriveURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\t"+req.URL.String())
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "\tDid you accidentally use an absolute endpoint URL rather than relative?")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	// server is a test HTTP server used to provide mock API responses.
	server := httptest.NewServer(apiHandler)

	// client is the GitHub client being tested and is configured to use test server.
	url, _ = url.Parse(server.URL + baseURLPath + "/")

	return url, mux, server.Close
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testHeader(t *testing.T, r *http.Request, header string, want string) {
	t.Helper()
	if got := r.Header.Get(header); got != want {
		t.Errorf("Header.Get(%q) returned %q, want %q", header, got, want)
	}
}

func testBody[T any](t *testing.T, r *http.Request, body T) {
	t.Helper()
	var got T
	err := json.NewDecoder(r.Body).Decode(&got)
	if err != nil {
		t.Errorf("Failed to decode request body: %v", err)
	}
	if !reflect.DeepEqual(got, body) {
		t.Errorf("Request body: %+v, want %+v", got, body)
	}
}

func readFile(t *testing.T, fileName string) []byte {
	testData, err := os.ReadFile("testdata/" + fileName)
	if err != nil {
		t.Errorf("readTestData failed: %v", err)
	}
	return testData
}

func getDataFromFile[T any](t *testing.T, fileName string) T {
	var data T
	err := json.Unmarshal(readFile(t, fileName), &data)
	if err != nil {
		t.Errorf("readTestData failed: %v", err)
	}
	return data
}

func getDataFromRequest[T any](t *testing.T, r *http.Request) T {
	var data T
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		t.Errorf("readTestData failed: %v", err)
	}
	return data
}

type fakeFile struct {
	readTimes int
}

func (f *fakeFile) Read(p []byte) (n int, err error) {
	nList := []int{11, fragmentSize - 10, fragmentSize, 1}
	errList := []error{nil, nil, nil, nil}

	if f.readTimes >= len(nList) {
		return 0, io.EOF
	}
	n = nList[f.readTimes]
	err = errList[f.readTimes]
	f.readTimes++
	return
}

func (f *fakeFile) Name() string {
	return "fake_file_name"
}

func (f *fakeFile) IsDir() bool {
	return false
}

func (f *fakeFile) Size() int64 {
	return fragmentSize*2 + 1
}

func (f *fakeFile) Finished() bool {
	return f.readTimes >= 4
}

func (f *fakeFile) ContentRange() string {
	template := "bytes %d-%d/%d"
	switch f.readTimes {
	case 1, 2:
		return fmt.Sprintf(template, 0, fragmentSize-1, f.Size())
	case 3:
		return fmt.Sprintf(template, fragmentSize, fragmentSize*2-1, f.Size())
	case 4:
		return fmt.Sprintf(template, fragmentSize*2, fragmentSize*2, f.Size())
	}
	return ""
}

func (f *fakeFile) ContentLength() string {
	switch f.readTimes {
	case 1, 2, 3:
		return strconv.Itoa(fragmentSize)
	case 4:
		return strconv.Itoa(1)
	}
	return ""
}

type fakeDir struct {
}

func (f *fakeDir) Name() string {
	return "fake_dir_name"
}

func (f *fakeDir) IsDir() bool {
	return true
}

func (f *fakeDir) Size() int64 {
	return 0
}

func (f *fakeDir) Read(p []byte) (n int, err error) {
	return 0, nil
}

type fakeEmptyFile struct {
}

func (f *fakeEmptyFile) Name() string {
	return "fake_empty_file_name"
}

func (f *fakeEmptyFile) IsDir() bool {
	return false
}

func (f *fakeEmptyFile) Size() int64 {
	return 0
}

func (f *fakeEmptyFile) Read(p []byte) (n int, err error) {
	return 0, nil
}

type fakeErrorFile struct {
}

func (f *fakeErrorFile) Name() string {
	return "fake_error_file_name"
}

func (f *fakeErrorFile) IsDir() bool {
	return false
}

func (f *fakeErrorFile) Size() int64 {
	return 1
}

func (f *fakeErrorFile) Read(p []byte) (n int, err error) {
	return 0, os.ErrNotExist
}
