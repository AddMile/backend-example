package httpclient

import "fmt"

// APIError represents an error returned from the API.
type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
}

// Error returns the error message.
func (err APIError) Error() string {
	return fmt.Sprintf("%d: %s", err.Status, err.Message)
}
