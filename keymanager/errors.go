package keymanager

import (
	"encoding/json"
	"log"

	"github.com/pkg/errors"
)

// HTTPRequestError represents an HTTP request error.
type HTTPRequestError struct {
	URL          string `json:"url,required"`
	StatusCode   int    `json:"status_code,omitempty"`
	ResponseBody []byte `json:"response_body,omitempty"`
	Message      string `json:"message,omitempty"`
}

// NewHTTPRequestError is the constructor of HTTPRequestError.
func NewHTTPRequestError(url string, statusCode int, responseBody []byte, message string) *HTTPRequestError {
	return &HTTPRequestError{
		URL:          url,
		StatusCode:   statusCode,
		ResponseBody: responseBody,
		Message:      message,
	}
}

// Error implements error interface
func (e *HTTPRequestError) Error() string {
	return e.String()
}

// String returns a readable string representation of a HTTPRequestError struct.
func (e *HTTPRequestError) String() string {
	if e == nil {
		return ""
	}

	data, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}

// GenericError represents the generic error of keymanager.
type GenericError struct {
	Err error `json:"err"`
}

// NewGenericError is the constructor of GenericError.
func NewGenericError(err error, desc string, args ...interface{}) *GenericError {
	return &GenericError{
		Err: errors.Wrapf(err, desc, args...),
	}
}

// NewGenericErrorWithMessage is the constructor of GenericError.
func NewGenericErrorWithMessage(msg string) *GenericError {
	return &GenericError{
		Err: errors.New(msg),
	}
}

// Error implements error interface.
func (e *GenericError) Error() string {
	return e.String()
}

// String implements fmt.Stringer interface.
func (e *GenericError) String() string {
	if e == nil {
		return ""
	}

	data, err := json.Marshal(e)
	if err != nil {
		log.Fatal(err)
	}
	return string(data)
}
