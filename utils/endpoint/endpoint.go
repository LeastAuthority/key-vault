package endpoint

const (
	// BasePath is the base path for all endpoints.
	BasePath = "/v1/ethereum"
)

// Build builds full path.
func Build(pattern string) string {
	return BasePath + "/" + pattern
}
