package endpoint

import "fmt"

const (
	// BasePath is the base path for all endpoints.
	BasePath = "/v1/ethereum"
)

// Build builds full path.
func Build(network, pattern string) string {
	if len(network) > 0 {
		return fmt.Sprintf("%s/%s/%s", BasePath, network, pattern)
	}

	return fmt.Sprintf("%s/%s", BasePath, pattern)
}
