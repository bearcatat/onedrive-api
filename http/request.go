package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Request interface {
	GetHttpRequest() (*http.Request, error)
}

type JsonRequest struct {
	method string
	body   interface{}
	url    *url.URL
}

func NewJsonRequest(method string, url *url.URL, body interface{}) Request {
	return &JsonRequest{
		method: method,
		url:    url,
		body:   body,
	}
}

func (r *JsonRequest) GetHttpRequest() (*http.Request, error) {
	if r.body != nil {
		return r.GetHttpRequestWithBody()
	}
	return http.NewRequest(r.method, r.url.String(), nil)
}

func (r *JsonRequest) GetHttpRequestWithBody() (*http.Request, error) {
	bodyReader, err := r.GetBodyReader()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(r.method, r.url.String(), bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

func (r *JsonRequest) GetBodyReader() (io.Reader, error) {
	jsonBody, err := json.Marshal(r.body)
	if err != nil {
		return nil, err
	}
	jsonBuffer := bytes.NewBuffer([]byte(jsonBody))
	return jsonBuffer, nil
}

type FileFragmentUploadReqeust struct {
	url   url.URL
	start int64
	total int64
	body  []byte
}

func NewFileFragmentUploadRequest(url url.URL, start, total int64, body []byte) Request {
	return &FileFragmentUploadReqeust{
		url:   url,
		start: start,
		total: total,
		body:  body,
	}
}

func (r *FileFragmentUploadReqeust) GetHttpRequest() (*http.Request, error) {
	reader := bytes.NewReader(r.body)
	req, err := http.NewRequest(http.MethodPut, r.url.String(), reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Length", strconv.FormatInt(reader.Size(), 10))
	end := r.start + reader.Size() - 1
	req.Header.Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", r.start, end, r.total))
	return req, nil
}
