package control

import (
	"fmt"
)

type ErrorInfo struct {
	Message    string              `json:"message"`
	Code       int                 `json:"code"`
	StatusCode int                 `json:"statusCode"`
	HRef       string              `json:"href"`
	Details    map[string][]string `json:"details"`
}

func (e ErrorInfo) Error() string {
	errorHref := e.HRef
	if e.HRef == "" && e.Code != 0 {
		errorHref = fmt.Sprintf("https://help.ably.io/error/%d", e.Code)
	}

	err := fmt.Sprintf("%s: code %d: status code: %d", e.Message, e.Code, e.StatusCode)
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
