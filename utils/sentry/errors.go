package sentry

// CustomComplexError implements Errors interface to provide custom complex error.
type CustomComplexError struct {
	Message  string
	MoreData map[string]string
}

// CustomComplexError implements Errors interface.
func (e CustomComplexError) Error() string {
	return e.Message
}
