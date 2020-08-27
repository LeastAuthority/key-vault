package httpex

import (
	"net/http"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// Default configuration of HTTP client.
const (
	attempts        = 3
	attemptsWaitMin = time.Second / 3
	attemptsWaitMax = time.Second
	clientTimeout   = time.Minute
)

// CreateClient creates a new HTTP client.
func CreateClient() *http.Client {
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = attempts
	retryClient.RetryWaitMin = attemptsWaitMin
	retryClient.RetryWaitMax = attemptsWaitMax

	client := retryClient.StandardClient()
	client.Timeout = clientTimeout

	return client
}
