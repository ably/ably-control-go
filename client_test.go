package control

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var token string
var url string = API_URL

func TestMain(m *testing.M) {
	token = os.Getenv("ABLY_ACCOUNT_TOKEN")

	if os.Getenv("ABLY_CONTROL_URL") != "" {
		url = os.Getenv("ABLY_CONTROL_URL")
	}

	os.Exit(m.Run())
}

func skipIntegrationTest(t *testing.T) {
	t.Helper()
	if os.Getenv("SKIP_INTEGRATION") != "" {
		t.Skip("SKIP_INTEGRATION is set, skipping integration test")
	}
}

func newTestApp(t *testing.T, client *Client) App {
	t.Helper()
	n := rand.Uint64()
	name := "test-" + fmt.Sprint(n)
	t.Logf("creating app with name: %s", name)
	app, err := client.CreateApp(&NewApp{
		Name:   name,
		Status: "enabled",
	})
	require.NoError(t, err)

	t.Cleanup(func() {
		client.DeleteApp(app.ID)
	})

	return app
}

func newTestClient(t *testing.T) (Client, Me) {
	t.Helper()
	skipIntegrationTest(t)
	client, me, err := NewClientWithURL(token, url)
	require.NoError(t, err)
	return client, me
}

// TestAblyAgent tests that client requests set the Ably-Agent HTTP header.
func TestAblyAgent(t *testing.T) {
	// start a test HTTP server which tracks the value of the Ably-Agent
	// HTTP header and returns an empty JSON object.
	var ablyAgent string
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ablyAgent = req.Header.Get("Ably-Agent")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", "2")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{}"))
	})
	srv := httptest.NewServer(handler)

	// initialise a client, which will make a request to /me
	client, _, err := NewClientWithURL("s3cr3t", srv.URL)
	assert.NoError(t, err)

	// check the Ably-Agent HTTP header was set
	assert.Equal(t, "ably-control-go/"+VERSION, ablyAgent)

	// add an extra Ably-Agent entry
	client.AppendAblyAgent("test", "1.2.3")

	// check requests now set the updated Ably-Agent HTTP header
	_, err = client.Me()
	assert.NoError(t, err)
	assert.Equal(t, "ably-control-go/"+VERSION+" test/1.2.3", ablyAgent)
}

// noRetryDelay returns a zero-delay backoff function for fast tests.
func noRetryDelay() Backoff {
	return func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		return 0
	}
}

// TestRetryOn5xxErrors tests that the client retries on 5xx server errors.
func TestRetryOn5xxErrors(t *testing.T) {
	var requestCount int32

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		if count < 3 {
			// First two requests return 503
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"message": "Service unavailable"}`))
			return
		}
		// Third request succeeds
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	client, _, err := NewClientWithURL("test-token", srv.URL,
		WithRetryMax(4),
		WithBackoff(noRetryDelay()),
	)

	assert.NoError(t, err)
	assert.Equal(t, int32(3), atomic.LoadInt32(&requestCount), "expected 3 requests (1 initial + 2 retries)")
	assert.NotNil(t, client)
}

// TestNoRetryOn4xxErrors tests that the client does NOT retry on 4xx client errors.
func TestNoRetryOn4xxErrors(t *testing.T) {
	var requestCount int32

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message": "Bad request", "code": 40000, "statusCode": 400}`))
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	_, _, err := NewClientWithURL("test-token", srv.URL,
		WithRetryMax(4),
		WithBackoff(noRetryDelay()),
	)

	assert.Error(t, err)
	assert.Equal(t, int32(1), atomic.LoadInt32(&requestCount), "expected exactly 1 request (no retries on 4xx)")

	errorInfo, ok := err.(ErrorInfo)
	assert.True(t, ok, "error should be ErrorInfo")
	assert.Equal(t, 400, errorInfo.StatusCode)
}

// TestRetryMaxExceeded tests that the client gives up after max retries.
func TestRetryMaxExceeded(t *testing.T) {
	var requestCount int32

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"message": "Service unavailable"}`))
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	_, _, err := NewClientWithURL("test-token", srv.URL,
		WithRetryMax(3),
		WithBackoff(noRetryDelay()),
	)

	assert.Error(t, err)
	// Initial request + 3 retries = 4 total requests
	assert.Equal(t, int32(4), atomic.LoadInt32(&requestCount), "expected 4 requests (1 initial + 3 retries)")
}

// TestRetryMaxOption tests that WithRetryMax correctly configures retry attempts.
func TestRetryMaxOption(t *testing.T) {
	var requestCount int32

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"message": "Service unavailable"}`))
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	// Test with RetryMax of 1
	_, _, err := NewClientWithURL("test-token", srv.URL,
		WithRetryMax(1),
		WithBackoff(noRetryDelay()),
	)

	assert.Error(t, err)
	// Initial request + 1 retry = 2 total requests
	assert.Equal(t, int32(2), atomic.LoadInt32(&requestCount), "expected 2 requests with RetryMax=1")
}

// TestBackoffConfiguration tests that backoff receives correct min/max configuration.
func TestBackoffConfiguration(t *testing.T) {
	var capturedMin, capturedMax time.Duration
	var capturedAttempts []int

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"message": "Service unavailable"}`))
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	expectedMin := 5 * time.Second
	expectedMax := 60 * time.Second

	customBackoff := func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		capturedMin = min
		capturedMax = max
		capturedAttempts = append(capturedAttempts, attemptNum)
		return 0 // Don't actually wait
	}

	_, _, _ = NewClientWithURL("test-token", srv.URL,
		WithRetryMax(2),
		WithRetryWaitMin(expectedMin),
		WithRetryWaitMax(expectedMax),
		WithBackoff(customBackoff),
	)

	assert.Equal(t, expectedMin, capturedMin, "backoff should receive configured min wait")
	assert.Equal(t, expectedMax, capturedMax, "backoff should receive configured max wait")
	assert.Equal(t, []int{0, 1}, capturedAttempts, "backoff should be called with attempt numbers 0 and 1")
}

// TestWithHTTPClient tests that a custom HTTP client can be provided.
func TestWithHTTPClient(t *testing.T) {
	customTransport := &http.Transport{
		MaxIdleConns: 100,
	}
	customHTTPClient := &http.Client{
		Transport: customTransport,
		Timeout:   30 * time.Second,
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	client, _, err := NewClientWithURL("test-token", srv.URL,
		WithHTTPClient(customHTTPClient),
	)

	assert.NoError(t, err)
	assert.NotNil(t, client)
	// Verify the custom HTTP client is being used by checking it was set
	assert.Equal(t, customHTTPClient, client.httpClient.HTTPClient)
}

// TestCheckRetryOption tests that a custom CheckRetry function can be provided.
func TestCheckRetryOption(t *testing.T) {
	var checkRetryCalled int32

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"message": "Service unavailable"}`))
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	customCheckRetry := func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		atomic.AddInt32(&checkRetryCalled, 1)
		// Never retry
		return false, nil
	}

	_, _, err := NewClientWithURL("test-token", srv.URL,
		WithRetryMax(3),
		WithCheckRetry(customCheckRetry),
		WithBackoff(noRetryDelay()),
	)

	assert.Error(t, err)
	// CheckRetry is called once and returns false, so only 1 request
	assert.Equal(t, int32(1), atomic.LoadInt32(&checkRetryCalled), "custom CheckRetry should be called")
}

// TestRetryOn429TooManyRequests tests that the client retries on 429 rate limit errors.
func TestRetryOn429TooManyRequests(t *testing.T) {
	var requestCount int32

	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		if count < 3 {
			// First two requests return 429
			w.Header().Set("Retry-After", "0")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"message": "Rate limit exceeded"}`))
			return
		}
		// Third request succeeds
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	client, _, err := NewClientWithURL("test-token", srv.URL,
		WithRetryMax(4),
		WithBackoff(noRetryDelay()),
	)

	assert.NoError(t, err)
	assert.Equal(t, int32(3), atomic.LoadInt32(&requestCount), "expected 3 requests (1 initial + 2 retries on 429)")
	assert.NotNil(t, client)
}

// TestRetryOnConnectionError tests that the client retries on connection errors.
// This simulates the status code 0 scenario that can occur with connection issues.
func TestRetryOnConnectionError(t *testing.T) {
	var requestCount int32

	// Start a server that we'll shut down to simulate connection errors
	handler := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		count := atomic.AddInt32(&requestCount, 1)
		if count < 3 {
			// Close connection abruptly to simulate connection error
			hj, ok := w.(http.Hijacker)
			if ok {
				conn, _, _ := hj.Hijack()
				conn.Close()
				return
			}
		}
		// Third request succeeds
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	})
	srv := httptest.NewServer(handler)
	defer srv.Close()

	client, _, err := NewClientWithURL("test-token", srv.URL,
		WithRetryMax(4),
		WithBackoff(noRetryDelay()),
	)

	assert.NoError(t, err)
	assert.Equal(t, int32(3), atomic.LoadInt32(&requestCount), "expected 3 requests (retries on connection error)")
	assert.NotNil(t, client)
}
