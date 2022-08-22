package control

import (
	"fmt"
)

// ErrorInfo represents an error type that the REST API may return.
type ErrorInfo struct {
	// The error message.
	Message string `json:"message"`
	// The Ably error code.
	Code int `json:"code"`
	// The HTTP status code returned.
	StatusCode int `json:"statusCode"`
	// The URL to documentation about the error code.
	HRef string `json:"href"`
	// Any additional details about the error message.
	Details map[string][]string `json:"details"`
	// The API path that resulted in this error.
	APIPath string
}

// ErrorInfo implements the Error interface.
func (e ErrorInfo) Error() string {
	errorHref := e.HRef
	if e.HRef == "" && e.Code != 0 {
		errorHref = fmt.Sprintf("https://help.ably.io/error/%d", e.Code)
	}

	err := fmt.Sprintf("%s: %s: code %d: status code: %d", e.APIPath, e.Message, e.Code, e.StatusCode)
	if errorHref != "" {
		err += fmt.Sprintf(" see: %s", errorHref)
	}
	if len(e.Details) != 0 {
		for k, v := range e.Details {
			err += fmt.Sprintf("\n  %s:", k)
			for _, str := range v {
				err += fmt.Sprintf("\n    %s", str)
			}
		}
	}
	return err
}
