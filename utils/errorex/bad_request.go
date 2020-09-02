package errorex

import (
	"net/http"

	"github.com/hashicorp/vault/sdk/logical"
)

// ErrBadRequest represents the bad request error
type ErrBadRequest struct {
	ErrorMsg string `json:"error_msg"`
}

// NewErrBadRequest is the constructor of ErrBadRequest
func NewErrBadRequest(errorMsg string) *ErrBadRequest {
	return &ErrBadRequest{
		ErrorMsg: errorMsg,
	}
}

// Error implements error interface
func (e *ErrBadRequest) Error() string {
	return e.ErrorMsg
}

// ToLogicalResponse converts error to logical response model
func (e *ErrBadRequest) ToLogicalResponse() (*logical.Response, error) {
	return logical.RespondWithStatusCode(&logical.Response{
		Data: map[string]interface{}{
			"message":     e.ErrorMsg,
			"status_code": http.StatusBadRequest,
		},
	}, nil, http.StatusBadRequest)
}
