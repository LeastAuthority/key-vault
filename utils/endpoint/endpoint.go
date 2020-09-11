package endpoint

import "fmt"

const (
	// BasePath is the base path for all endpoints.
	BasePath = "/v1/ethereum"
)

// Build builds full path.
func Build(network, pattern string) string {
	return fmt.Sprintf("%s/%s/%s", BasePath, network, pattern)
}
