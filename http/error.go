package http

import (
	"fmt"
)

// ErrorResponse represents the error response returned by OneDrive drive API.
type ErrorResponse struct {
	Error *Error `json:"error"`
}

func (r *ErrorResponse) GetError() error {
	if r == nil || r.Error == nil {
		return nil
	}
	if r.Error.InnerError != nil {
		return fmt.Errorf("%s-%s (%s)", r.Error.Code, r.Error.Message, r.Error.InnerError.Date)
	}
	return fmt.Errorf("%s-%s", r.Error.Code, r.Error.Message)
}

// Error represents the error in the response returned by OneDrive drive API.
type Error struct {
	Code             string      `json:"code"`
	Message          string      `json:"message"`
	LocalizedMessage string      `json:"localizedMessage"`
	InnerError       *InnerError `json:"innerError"`
}

// InnerError represents the error details in the error returned by OneDrive drive API.
type InnerError struct {
	Date            string `json:"date"`
	RequestId       string `json:"request-id"`
	ClientRequestId string `json:"client-request-id"`
}
