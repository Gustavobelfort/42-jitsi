package intra

import (
	"fmt"
	"net/http"
)

func validateResponse(resp *http.Response) error {
	if 200 > resp.StatusCode || resp.StatusCode > 299 {
		return &HTTPError{Response: resp}
	}
	return nil
}

// HTTPError wraps a bad http response into a golang error.
type HTTPError struct {
	Response *http.Response
}

// Error returns the formatted http error.
func (err *HTTPError) Error() string {
	return fmt.Sprintf(
		"%s %s: %d: %s",
		err.Response.Request.Method,
		err.Response.Request.URL.String(),
		err.Response.StatusCode,
		err.Response.Status)
}
