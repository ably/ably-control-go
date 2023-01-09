package control

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var token string
var url string = API_URL
var apps []string

func TestMain(m *testing.M) {
	token = os.Getenv("ABLY_ACCOUNT_TOKEN")
	rand.Seed(time.Now().UnixNano())

	if token == "" {
		panic("ABLY_ACCOUNT_TOKEN not set")
	}

	if os.Getenv("ABLY_CONTROL_URL") != "" {
		url = os.Getenv("ABLY_CONTROL_URL")
	}

	// Attempt to clean up apps if anything went wrong
	client, _, err := NewClientWithURL(token, url)
	if err == nil {
		for _, v := range apps {
			_ = client.DeleteApp(v)
		}
	}

	code := m.Run()
	os.Exit(code)
}

func newTestApp(t *testing.T, client *Client) App {
	n := rand.Uint64()
	name := "test-" + fmt.Sprint(n)
	t.Logf("creating app with name: %s", name)
	app := NewApp{
		Name:   name,
		Status: "enabled",
		//TLSOnly:                false,
		FcmKey:                 "",
		ApnsCertificate:        "",
		ApnsPrivateKey:         "",
		ApnsUseSandboxEndpoint: false,
	}
	app_ret, err := client.CreateApp(&app)

	assert.NoError(t, err)
	apps = append(apps, app.ID)

	return app_ret
}

func newTestClient(t *testing.T) (Client, Me) {
	client, me, err := NewClientWithURL(token, url)
	assert.NoError(t, err)
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
