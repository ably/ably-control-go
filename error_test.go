package control

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	_, _, err := NewClientWithURL("", url)
	errorInfo := err.(ErrorInfo)

	expected := ErrorInfo{
		Message:    "Access denied",
		Code:       40100,
		StatusCode: 401,
		HRef:       "https://help.ably.io/error/40100",
		APIPath:    "/me",
	}

	assert.Equal(t, expected, errorInfo)
}

func TestErrorInfoWithInterfaceDetails(t *testing.T) {
	e := ErrorInfo{
		Message:    "Validation failed",
		Code:       40000,
		StatusCode: 400,
		Details: map[string]any{
			"fields": []any{"name is required", "status is invalid"},
			"count":  42.0,
		},
		APIPath: "/apps",
	}

	errStr := e.Error()
	assert.Contains(t, errStr, "Validation failed")
	assert.Contains(t, errStr, "name is required")
	assert.Contains(t, errStr, "status is invalid")
	assert.Contains(t, errStr, "42")
}

func TestErrorInfoUnmarshalDetails(t *testing.T) {
	data := `{
		"message": "Bad request",
		"code": 40000,
		"statusCode": 400,
		"href": "",
		"details": {
			"errors": ["field1 invalid", "field2 missing"],
			"severity": "high"
		}
	}`

	var e ErrorInfo
	err := json.Unmarshal([]byte(data), &e)
	assert.NoError(t, err)
	assert.NotNil(t, e.Details)

	errors, ok := e.Details["errors"].([]any)
	assert.True(t, ok)
	assert.Len(t, errors, 2)
	assert.Equal(t, "field1 invalid", errors[0])

	severity, ok := e.Details["severity"].(string)
	assert.True(t, ok)
	assert.Equal(t, "high", severity)
}
