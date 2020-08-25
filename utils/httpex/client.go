package httpex

import (
	"net"
	"net/http"
	"time"
)

// Default configuration of HTTP client.
const (
	dealTimeout         = time.Minute
	tlsHandshakeTimeout = time.Minute
	clientTimeout       = time.Minute
)

// CreateClient creates a new HTTP client.
func CreateClient() *http.Client {
	netTransport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: dealTimeout,
		}).Dial,
		TLSHandshakeTimeout: tlsHandshakeTimeout,
	}

	return &http.Client{
		Timeout:   clientTimeout,
		Transport: netTransport,
	}
}
