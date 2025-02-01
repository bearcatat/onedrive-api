package onedrive

import (
	"io"
	"net/url"

	"github.com/bearcatat/onedrive-api/http"
	"github.com/bearcatat/onedrive-api/resources"
)

const (
	fragmentSize = 10 * 1024 * 1024
)

type File interface {
	io.Reader
	Name() string
	IsDir() bool
	Size() int64
}
type fileForUpload struct {
	f            File
	size         int64
	uploadedSize int64
	url          *url.URL
}

func newFileForUpload(f File, uploadSession *resources.UploadSession) *fileForUpload {
	file := &fileForUpload{
		f:            f,
		uploadedSize: 0,
	}
	file.url, _ = url.Parse(uploadSession.UploadURL)
	file.size = f.Size()
	return file
}

func (f *fileForUpload) next() bool {
	return f.uploadedSize < f.size
}

func (f *fileForUpload) getNextRequest() (http.Request, error) {
	fragment, err := f.readFragment()
	if err != nil {
		return nil, err
	}
	req := http.NewFileFragmentUploadRequest(*f.url, f.uploadedSize, f.f.Size(), fragment)
	f.uploadedSize += int64(len(fragment))
	return req, nil
}

func (f *fileForUpload) readFragment() ([]byte, error) {
	buffer := make([]byte, fragmentSize)
	n := 0
	for n < fragmentSize {
		sn, err := f.f.Read(buffer[n:])
		if err != nil {
			if err == io.EOF {
				n += sn
				break
			}
			return nil, err
		}
		n += sn
	}
	if n < fragmentSize {
		bufferLast := buffer[:n]
		buffer = bufferLast
	}
	return buffer, nil
}
