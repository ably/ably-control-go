package control

import (
	"fmt"
	"math/rand"
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
