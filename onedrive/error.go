package onedrive

import "errors"

var (
	ErrNotFile     = errors.New("not a file")
	ErrEmptyFile   = errors.New("empty file")
	ErrNotFinished = errors.New("not finished")
)
