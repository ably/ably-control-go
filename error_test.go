package control

import (
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
	}

	assert.Equal(t, expected, errorInfo)
}
